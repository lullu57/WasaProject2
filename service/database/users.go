package database

//All User related methods are defined here

import (
	"crypto/rand"
	"database/sql"
	"fmt"
)

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}

// checkUserIDExists now returns an error as well
func (db *appdbimpl) checkUserIDExists(userID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", userID).Scan(&exists)
	return exists, err // return the error
}

// generateUniqueID now has a receiver and returns errors
func (db *appdbimpl) generateUniqueID() (string, error) {
	for {
		userID, err := generateRandomString(10)
		if err != nil {
			return "", err
		}

		exists, err := db.checkUserIDExists(userID) // use the method of appdbimpl
		if err != nil {
			return "", err // return the error
		}

		if !exists {
			return userID, nil
		}
	}
}

// AddUser uses the generateUniqueID method of appdbimpl
func (db *appdbimpl) AddUser(user *User) error {
	userID, err := db.generateUniqueID() // use the method of appdbimpl
	if err != nil {
		return fmt.Errorf("failed to generate user ID: %w", err)
	}
	user.ID = userID

	// Prepare SQL statement for inserting a new user
	stmt, err := db.c.Prepare("INSERT INTO users (user_id, username) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(user.ID, user.Username)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	// User ID is already set in the user struct, no need to return it
	return nil
}

func (db *appdbimpl) SetUsername(currentUsername, newUsername string) error {
	stmt, err := db.c.Prepare("UPDATE users SET username = ? WHERE username = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUsername, currentUsername)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (db *appdbimpl) GetUserProfile(username string) (*User, error) {
	var user User

	// Fetch basic user info
	err := db.c.QueryRow("SELECT user_id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Fetch followers
	rows, err := db.c.Query("SELECT follower_id FROM followers WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("failed to read follower row: %w", err)
		}
		user.Followers = append(user.Followers, followerID)
	}

	// Fetch following
	rows, err = db.c.Query("SELECT user_id FROM followers WHERE follower_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch following: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followingID string
		if err := rows.Scan(&followingID); err != nil {
			return nil, fmt.Errorf("failed to read following row: %w", err)
		}
		user.Following = append(user.Following, followingID)
	}

	// Fetch photos
	rows, err = db.c.Query("SELECT photo_id FROM user_photos WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch photos: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var photoID string
		if err := rows.Scan(&photoID); err != nil {
			return nil, fmt.Errorf("failed to read photo row: %w", err)
		}
		user.Photos = append(user.Photos, photoID)
	}

	return &user, nil
}

func (db *appdbimpl) FollowUser(followerID, followedID string) error {
	_, err := db.c.Exec(`INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`, followedID, followerID)
	if err != nil {
		return fmt.Errorf("error following user: %w", err)
	}
	return nil
}

func (db *appdbimpl) UnfollowUser(followerID, followedID string) error {
	_, err := db.c.Exec(`DELETE FROM followers WHERE user_id = ? AND follower_id = ?`, followedID, followerID)
	if err != nil {
		return fmt.Errorf("error unfollowing user: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetUserIDByUsername(username string) (string, error) {
	var userID string
	err := db.c.QueryRow("SELECT user_id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("query error: %w", err)
	}
	return userID, nil
}
