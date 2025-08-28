package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fact-check/internal/config"
	"fact-check/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	config       *config.Config
	db           *sql.DB
	logger       *logrus.Logger
	oauth2Config *oauth2.Config
}

func NewAuthService(cfg *config.Config, db *sql.DB, logger *logrus.Logger) *AuthService {
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		config:       cfg,
		db:           db,
		logger:       logger,
		oauth2Config: oauth2Config,
	}
}

func (s *AuthService) GetLoginURL() string {
	return s.oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleCallback(code string) (*models.AuthCallbackResponse, error) {
	s.logger.Infof("Starting OAuth callback with code: %s", code[:10]+"...")

	// Exchange code for token
	token, err := s.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Errorf("OAuth exchange failed: %v", err)
		s.logger.Errorf("OAuth config - ClientID: %s, RedirectURL: %s", s.oauth2Config.ClientID, s.oauth2Config.RedirectURL)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	s.logger.Infof("OAuth exchange successful, got access token: %s", token.AccessToken[:10]+"...")

	// Get user info from Google
	userInfo, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		s.logger.Errorf("Failed to get Google user info: %v", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	s.logger.Infof("Got user info for: %s (%s)", userInfo.Email, userInfo.Name)

	// Find or create user in database
	user, err := s.findOrCreateUser(userInfo)
	if err != nil {
		s.logger.Errorf("Failed to find or create user: %v", err)
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token
	jwtToken, err := s.generateJWT(user.ID.String())
	if err != nil {
		s.logger.Errorf("Failed to generate JWT: %v", err)
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	s.logger.Infof("OAuth callback completed successfully for user: %s", user.Email)

	return &models.AuthCallbackResponse{
		Token: jwtToken,
		User:  *user,
	}, nil
}

func (s *AuthService) getGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *AuthService) findOrCreateUser(googleUser *models.GoogleUserInfo) (*models.User, error) {
	// Try to find existing user
	var user models.User
	query := `SELECT id, google_id, email, name, picture, created_at, updated_at 
			  FROM users WHERE google_id = $1`

	err := s.db.QueryRow(query, googleUser.ID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		// User exists, update if needed
		if user.Name != googleUser.Name || user.Picture != googleUser.Picture {
			updateQuery := `UPDATE users SET name = $1, picture = $2, updated_at = CURRENT_TIMESTAMP 
						   WHERE id = $3 RETURNING updated_at`
			s.db.QueryRow(updateQuery, googleUser.Name, googleUser.Picture, user.ID).Scan(&user.UpdatedAt)
		}
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new user
	user = models.User{
		ID:        uuid.New(),
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		Picture:   googleUser.Picture,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertQuery := `INSERT INTO users (id, google_id, email, name, picture, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = s.db.Exec(insertQuery, user.ID, user.GoogleID, user.Email, user.Name,
		user.Picture, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *AuthService) generateJWT(userID string) (string, error) {
	s.logger.Infof("Generating JWT for user: %s", userID)
	s.logger.Infof("Using JWT secret: %s...", s.config.JWTSecret[:10]+"...")

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		s.logger.Errorf("Failed to sign JWT: %v", err)
		return "", err
	}

	s.logger.Infof("JWT generated successfully: %s...", tokenString[:20]+"...")
	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	s.logger.Infof("Validating token: %s...", tokenString[:10]+"...")
	s.logger.Infof("Using JWT secret: %s...", s.config.JWTSecret[:10]+"...")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		s.logger.Errorf("Token parsing failed: %v", err)
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		s.logger.Errorf("Token validation failed: token not valid")
		return "", fmt.Errorf("invalid token")
	}

	// Extract claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if subject, exists := claims["sub"]; exists {
			s.logger.Infof("Token validation successful for user: %s", subject)
			return subject.(string), nil
		} else {
			s.logger.Errorf("Token validation failed: subject claim not found")
			return "", fmt.Errorf("subject claim not found")
		}
	}

	s.logger.Errorf("Token validation failed: claims not found")
	return "", fmt.Errorf("invalid token")
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, google_id, email, name, picture, created_at, updated_at 
			  FROM users WHERE id = $1`

	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
