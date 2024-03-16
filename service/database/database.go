/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

type User struct {
	ID        string   `json:"userId"` // Unique identifier
	Username  string   `json:"username"`
	Followers []string // IDs of followers
	Following []string // IDs of users being followed
	Photos    []string // IDs of photos uploaded by the user
}

type Comment struct {
	ID        string    `json:"commentId"` // Unique identifier
	UserID    string    `json:"userId"`    // ID of the user who commented
	PhotoID   string    `json:"photoId"`   // ID of the photo being commented on
	Content   string    `json:"content"`   // The comment itself
	Timestamp time.Time `json:"timestamp"` // Timestamp of when the comment was made
}

type Like struct {
	UserID    string    `json:"userId"`    // ID of the user who liked the photo
	PhotoID   string    `json:"photoId"`   // ID of the photo being liked
	Timestamp time.Time `json:"timestamp"` // Timestamp of when the like was made
}

type Photo struct {
	ID         string    `json:"photoId"`    // Unique identifier
	UserID     string    `json:"userId"`     // ID of the user who uploaded the photo
	URL        string    `json:"url"`        // URL of the photo
	UploadTime time.Time `json:"uploadTime"` // Timestamp of when the photo was uploaded
	Likes      []Like
	Comments   []Comment
}

type Ban struct {
	ID         string    `json:"banId"`      // Unique identifier
	BannedBy   string    `json:"bannedBy"`   // ID of the user who banned the other user
	BannedUser string    `json:"bannedUser"` // ID of the user who was banned'
	Time       time.Time `json:"time"`       // Timestamp of when the ban was made
}
