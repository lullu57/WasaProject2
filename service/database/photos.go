package database

import (
	"fmt"
)

// AddPhoto stores metadata about a photo in the database.
func (db *appdbimpl) AddPhoto(photo Photo) error {
	stmt, err := db.c.Prepare("INSERT INTO new_photos (photo_id, user_id, image_data, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare the photo insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(photo.ID, photo.UserID, photo.ImageData, photo.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to execute the photo insert statement: %w", err)
	}

	return nil
}

// function to get all photos
func (db *appdbimpl) GetPhotos() ([]Photo, error) {
	rows, err := db.c.Query("SELECT * FROM new_photos")
	if err != nil {
		return nil, fmt.Errorf("failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.ID, &photo.UserID, &photo.ImageData, &photo.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (db *appdbimpl) DeletePhoto(photoID string) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// Delete comments
	_, err = tx.Exec("DELETE FROM comments WHERE photo_id = ?", photoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete likes
	_, err = tx.Exec("DELETE FROM likes WHERE photo_id = ?", photoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the photo
	_, err = tx.Exec("DELETE FROM new_photos WHERE photo_id = ?", photoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *appdbimpl) GetMyStream(userID string) ([]Photo, error) {
	var photos []Photo
	query := `SELECT p.* FROM new_photos p
              JOIN followers f ON p.user_id = f.user_id
              WHERE f.follower_id = ?`
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		if err := rows.Scan(&photo.ID, &photo.UserID, &photo.ImageData, &photo.Timestamp); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, nil
}
