package database

import (
	"crypto/rand"
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
