package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Retry connection with exponential backoff
	maxRetries := 30
	retryDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("Database connection attempt %d failed: %v", i+1, err)
			if i < maxRetries-1 {
				log.Printf("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				retryDelay = time.Duration(float64(retryDelay) * 1.5) // Exponential backoff
				if retryDelay > 10*time.Second {
					retryDelay = 10 * time.Second // Cap at 10 seconds
				}
				continue
			}
			return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
		}
		log.Printf("Database connection established successfully")
		break
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		google_id VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		picture VARCHAR(500),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	// Create news table
	createNewsTable := `
	CREATE TABLE IF NOT EXISTS news (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		link VARCHAR(500),
		photo_url VARCHAR(500),
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'true', 'false')),
		explanation TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_news_user_id ON news(user_id);
	CREATE INDEX IF NOT EXISTS idx_news_status ON news(status);
	CREATE INDEX IF NOT EXISTS idx_news_created_at ON news(created_at);
	CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`

	// Execute migrations
	migrations := []string{createUsersTable, createNewsTable, createIndexes}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
