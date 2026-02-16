package profile

import (
	"database/sql"
	"fmt"
	"mcli/internal/types"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const defaultDBPath = "data/mcli.db"

// Store wraps a SQLite connection for profile persistence
type Store struct {
	db *sql.DB
}

// OpenStore opens (or creates) the SQLite database and runs migrations
func OpenStore() (*Store, error) {
	dbPath := defaultDBPath
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	migrations := `
	CREATE TABLE IF NOT EXISTS profiles (
		user_id    TEXT PRIMARY KEY,
		location   TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS bookmarks (
		user_id  TEXT NOT NULL,
		event_id TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, event_id),
		FOREIGN KEY (user_id) REFERENCES profiles(user_id)
	);

	CREATE TABLE IF NOT EXISTS read_events (
		user_id  TEXT NOT NULL,
		event_id TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, event_id),
		FOREIGN KEY (user_id) REFERENCES profiles(user_id)
	);

	CREATE TABLE IF NOT EXISTS filters (
		user_id TEXT NOT NULL,
		name    TEXT NOT NULL,
		value   TEXT NOT NULL,
		PRIMARY KEY (user_id, name),
		FOREIGN KEY (user_id) REFERENCES profiles(user_id)
	);
	`
	_, err := s.db.Exec(migrations)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// Load retrieves a user profile from the database, creating one if it doesn't exist
func (s *Store) Load(userID string) (*UserProfile, error) {
	p := New(userID)

	// Try to load existing profile
	var location string
	var createdAt, updatedAt time.Time
	err := s.db.QueryRow(
		"SELECT location, created_at, updated_at FROM profiles WHERE user_id = ?", userID,
	).Scan(&location, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		// Insert new profile
		now := time.Now()
		_, err = s.db.Exec(
			"INSERT INTO profiles (user_id, location, created_at, updated_at) VALUES (?, '', ?, ?)",
			userID, now, now,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create profile: %w", err)
		}
		return p, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load profile: %w", err)
	}

	p.Location = location
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt

	// Load bookmarks
	rows, err := s.db.Query("SELECT event_id FROM bookmarks WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load bookmarks: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var eid string
		if err := rows.Scan(&eid); err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}
		p.Bookmarks = append(p.Bookmarks, types.EventId(eid))
	}

	// Load read events
	rows, err = s.db.Query("SELECT event_id FROM read_events WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load read events: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var eid string
		if err := rows.Scan(&eid); err != nil {
			return nil, fmt.Errorf("failed to scan read event: %w", err)
		}
		p.ReadEvents = append(p.ReadEvents, types.EventId(eid))
	}

	// Load filters
	rows, err = s.db.Query("SELECT name, value FROM filters WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load filters: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, fmt.Errorf("failed to scan filter: %w", err)
		}
		p.Filters[name] = value
	}

	return p, nil
}

// Save persists the profile and all related data to the database
func (s *Store) Save(p *UserProfile) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	p.UpdatedAt = now

	// Upsert profile
	_, err = tx.Exec(`
		INSERT INTO profiles (user_id, location, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET location = ?, updated_at = ?`,
		p.UserID, p.Location, p.CreatedAt, now,
		p.Location, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	// Replace bookmarks
	if _, err := tx.Exec("DELETE FROM bookmarks WHERE user_id = ?", p.UserID); err != nil {
		return fmt.Errorf("failed to clear bookmarks: %w", err)
	}
	for _, eid := range p.Bookmarks {
		if _, err := tx.Exec("INSERT INTO bookmarks (user_id, event_id) VALUES (?, ?)", p.UserID, string(eid)); err != nil {
			return fmt.Errorf("failed to save bookmark: %w", err)
		}
	}

	// Replace read events
	if _, err := tx.Exec("DELETE FROM read_events WHERE user_id = ?", p.UserID); err != nil {
		return fmt.Errorf("failed to clear read events: %w", err)
	}
	for _, eid := range p.ReadEvents {
		if _, err := tx.Exec("INSERT INTO read_events (user_id, event_id) VALUES (?, ?)", p.UserID, string(eid)); err != nil {
			return fmt.Errorf("failed to save read event: %w", err)
		}
	}

	// Replace filters
	if _, err := tx.Exec("DELETE FROM filters WHERE user_id = ?", p.UserID); err != nil {
		return fmt.Errorf("failed to clear filters: %w", err)
	}
	for name, value := range p.Filters {
		if _, err := tx.Exec("INSERT INTO filters (user_id, name, value) VALUES (?, ?, ?)", p.UserID, name, value); err != nil {
			return fmt.Errorf("failed to save filter: %w", err)
		}
	}

	return tx.Commit()
}

// SaveLocation updates just the location field
func (s *Store) SaveLocation(userID, location string) error {
	_, err := s.db.Exec(
		"UPDATE profiles SET location = ?, updated_at = ? WHERE user_id = ?",
		location, time.Now(), userID,
	)
	return err
}

// AddBookmark adds a single bookmark
func (s *Store) AddBookmark(userID string, eventID types.EventId) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO bookmarks (user_id, event_id) VALUES (?, ?)",
		userID, string(eventID),
	)
	return err
}

// RemoveBookmark removes a single bookmark
func (s *Store) RemoveBookmark(userID string, eventID types.EventId) error {
	_, err := s.db.Exec(
		"DELETE FROM bookmarks WHERE user_id = ? AND event_id = ?",
		userID, string(eventID),
	)
	return err
}

// AddReadEvent marks a single event as read
func (s *Store) AddReadEvent(userID string, eventID types.EventId) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO read_events (user_id, event_id) VALUES (?, ?)",
		userID, string(eventID),
	)
	return err
}
