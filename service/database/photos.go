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

func (db *appdbimpl) GetPhoto(photoId string) (*PhotoDetail, error) {
	var photo PhotoDetail

	// First, fetch the basic photo details and count of likes
	err := db.c.QueryRow(`
    SELECT p.photo_id, p.user_id, u.username, p.image_data, p.timestamp,
           (SELECT COUNT(*) FROM likes WHERE photo_id = p.photo_id) AS likes_count
    FROM new_photos p
    JOIN users u ON p.user_id = u.user_id
    WHERE p.photo_id = ?`, photoId).Scan(
		&photo.PhotoID, &photo.UserID, &photo.Username, &photo.ImageData, &photo.Timestamp, &photo.LikesCount,
	)
	if err != nil {
		return nil, err
	}

	// Query for comments related to the photo
	commentsQuery := `
    SELECT c.comment_id, c.user_id, u.username, c.content, c.timestamp
    FROM comments c
    JOIN users u ON u.user_id = c.user_id
    WHERE c.photo_id = ?
    ORDER BY c.timestamp DESC
    `
	rows, err := db.c.Query(commentsQuery, photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photo.Comments = []Comment{} // Initialize the slice to store comments

	// Iterate over the results and populate the comments slice
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PhotoID, &comment.Content, &comment.Timestamp); err != nil {
			return nil, err
		}
		photo.Comments = append(photo.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &photo, nil
}
