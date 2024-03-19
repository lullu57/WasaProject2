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
	"time"
)

type Error struct {
	Error string `json:"error" db:"error"`
}

type User struct {
	ID        string   `json:"userId" db:"user_id"` // Unique identifier
	Username  string   `json:"username" db:"username"`
	Followers []string `json:"followers"` // IDs of followers (handled separately in relational mapping)
	Following []string `json:"following"` // IDs of users being followed (handled separately in relational mapping)
	Photos    []string `json:"photos"`    // IDs of photos uploaded by the user (handled separately in relational mapping)
}

// New Struct for handling followers relationship
type Follower struct {
	UserID     string `json:"userId" db:"user_id"`
	FollowerID string `json:"followerId" db:"follower_id"`
}

// New Struct for handling user and photo relationship
type UserPhoto struct {
	UserID  string `json:"userId" db:"user_id"`
	PhotoID string `json:"photoId" db:"photo_id"`
}

type Comment struct {
	ID        string    `json:"commentId" db:"comment_id"` // Unique identifier
	UserID    string    `json:"userId" db:"user_id"`       // ID of the user who commented
	PhotoID   string    `json:"photoId" db:"photo_id"`     // ID of the photo being commented on
	Content   string    `json:"content" db:"content"`      // The comment itself
	Timestamp time.Time `json:"timestamp" db:"timestamp"`  // Timestamp of when the comment was made
}

type Like struct {
	UserID    string    `json:"userId" db:"user_id"`      // ID of the user who liked the photo
	PhotoID   string    `json:"photoId" db:"photo_id"`    // ID of the photo being liked
	Timestamp time.Time `json:"timestamp" db:"timestamp"` // Timestamp of when the like was made
}

type Photo struct {
	ID        string    `json:"photoId" db:"photo_id"`    // Unique identifier
	UserID    string    `json:"userId" db:"user_id"`      // ID of the user who uploaded the photo
	URL       string    `json:"url" db:"url"`             // URL of the photo
	Timestamp time.Time `json:"timestamp" db:"timestamp"` // Timestamp of when the photo was uploaded
	Likes     []Like    `json:"likes"`                    // Note: This requires a relational mapping and isn't directly mapped to a single column
	Comments  []Comment `json:"comments"`                 // Note: This requires a relational mapping and isn't directly mapped to a single column
}

type Ban struct {
	ID         string    `json:"banId" db:"ban_id"`           // Unique identifier
	BannedBy   string    `json:"bannedBy" db:"banned_by"`     // ID of the user who banned the other user
	BannedUser string    `json:"bannedUser" db:"banned_user"` // ID of the user who was banned
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`    // Timestamp of when the ban was made
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	AddUser(user *User) error
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

	// Error table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS errors (
        error TEXT
    );`)
	if err != nil {
		return nil, err
	}

	// Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        user_id TEXT PRIMARY KEY,
        username TEXT UNIQUE NOT NULL
    );`)
	if err != nil {
		return nil, err
	}

	// Followers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS followers (
        user_id TEXT NOT NULL,
        follower_id TEXT NOT NULL,
        PRIMARY KEY (user_id, follower_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (follower_id) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	// User Photos table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_photos (
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        PRIMARY KEY (user_id, photo_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Comment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
        comment_id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        content TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Like table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS likes (
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        PRIMARY KEY (user_id, photo_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Photo table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS photos (
        photo_id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        url TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Ban table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bans (
        ban_id TEXT PRIMARY KEY,
        banned_by TEXT NOT NULL,
        banned_user TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (banned_by) REFERENCES users(user_id),
        FOREIGN KEY (banned_user) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

/*doLogin: Handle user login. If the user does not exist, it will create and log the user in, returning a unique user identifier. This method corresponds to the simplified login API spec provided.

setMyUserName: Allow a user to change their username. This requires updating the user's username in the database and handling any necessary validations.

uploadPhoto: Enable users to upload photos. This method involves storing the photo information in the database and handling the photo data itself, likely storing it in a file system or a blob storage service.

followUser: Allow a user to follow another user. This method requires updating the followers/following relationships in the database.

unfollowUser: Allow a user to unfollow another user. Similar to followUser, but in reverse.

banUser: Enable a user to ban another user, preventing the banned user from seeing any information about the banner.

unbanUser: Allow a user to remove a ban on another user.

getUserProfile: Retrieve and display a user's profile information, including their uploaded photos, the count of photos uploaded, and their followers/following lists.

getMyStream: Fetch and display a stream of photos from users that the current user follows, in reverse chronological order, including like and comment counts.

likePhoto: Allow a user to place a "like" on a photo.

unlikePhoto: Allow a user to remove their "like" from a photo.

commentPhoto: Enable users to add comments to a photo.

uncommentPhoto: Allow users to delete their comments from a photo.

deletePhoto: Allow users to remove a photo they have uploaded. This will also require removing all likes and comments associated with that photo.

*/
