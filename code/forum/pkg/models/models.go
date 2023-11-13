package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrDuplicateUsername = errors.New("models: duplicate username")

type Post struct {
	ID      int
	Author  string
	Title   string
	Body    string
	Image   string
	Votes   int
	Created time.Time
}

type Category struct {
	ID   int
	Name string
}

type PostCategory struct {
	CategoryID int
	Name       string
}

type Comment struct {
	ID       int
	PostID   int
	Body     string
	UserID   int
	Created  time.Time
	Votes    int
	Username string
}

type PostVote struct {
	UserID int
	PostID int
	Vote   int
}

type CommentVote struct {
	UserID    int
	CommentID int
	Vote      int
}

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type Session struct {
	UserID  int
	Token   string
	Created time.Time
	Expires time.Time
}
