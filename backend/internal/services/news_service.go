package services

import (
	"database/sql"
	"fmt"
	"time"

	"fact-check/internal/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NewsService struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewNewsService(cfg interface{}, db *sql.DB, logger *logrus.Logger) *NewsService {
	return &NewsService{
		db:     db,
		logger: logger,
	}
}

func (s *NewsService) SubmitNews(userID string, submission *models.NewsSubmission) (*models.News, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	news := &models.News{
		ID:        uuid.New(),
		UserID:    userUUID,
		Content:   submission.Content,
		Link:      submission.Link,
		PhotoURL:  submission.PhotoURL,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO news (id, user_id, content, link, photo_url, status, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = s.db.Exec(query, news.ID, news.UserID, news.Content, news.Link,
		news.PhotoURL, news.Status, news.CreatedAt, news.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert news: %w", err)
	}

	s.logger.Infof("News submitted successfully: %s", news.ID)
	return news, nil
}

func (s *NewsService) GetNewsByID(newsID string) (*models.News, error) {
	newsUUID, err := uuid.Parse(newsID)
	if err != nil {
		return nil, fmt.Errorf("invalid news ID: %w", err)
	}

	var news models.News
	query := `SELECT id, user_id, content, link, photo_url, status, explanation, created_at, updated_at 
			  FROM news WHERE id = $1`

	err = s.db.QueryRow(query, newsUUID).Scan(
		&news.ID, &news.UserID, &news.Content, &news.Link, &news.PhotoURL,
		&news.Status, &news.Explanation, &news.CreatedAt, &news.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("news not found")
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return &news, nil
}

func (s *NewsService) GetUserNews(userID string) ([]*models.News, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	query := `SELECT id, user_id, content, link, photo_url, status, explanation, created_at, updated_at 
			  FROM news WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user news: %w", err)
	}
	defer rows.Close()

	var newsList []*models.News
	for rows.Next() {
		var news models.News
		err := rows.Scan(
			&news.ID, &news.UserID, &news.Content, &news.Link, &news.PhotoURL,
			&news.Status, &news.Explanation, &news.CreatedAt, &news.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news row: %w", err)
		}
		newsList = append(newsList, &news)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over news rows: %w", err)
	}

	return newsList, nil
}

func (s *NewsService) UpdateNewsStatus(newsID string, status string, explanation string) error {
	newsUUID, err := uuid.Parse(newsID)
	if err != nil {
		return fmt.Errorf("invalid news ID: %w", err)
	}

	query := `UPDATE news SET status = $1, explanation = $2, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $3`

	result, err := s.db.Exec(query, status, explanation, newsUUID)
	if err != nil {
		return fmt.Errorf("failed to update news status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("news not found")
	}

	s.logger.Infof("News status updated successfully: %s -> %s", newsID, status)
	return nil
}
