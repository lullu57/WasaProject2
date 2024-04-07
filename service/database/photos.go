package database

import (
	"fmt"
)

// AddPhoto stores metadata about a photo in the database.
func (db *appdbimpl) AddPhoto(photo *Photo) error {
	stmt, err := db.c.Prepare("INSERT INTO photos (photo_id, user_id, url, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(photo.ID, photo.UserID, photo.URL, photo.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
