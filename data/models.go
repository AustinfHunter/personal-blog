package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID             int    `db:"ID" json:"id"`
	FirstName      string `db:"First_Name" json:"firstName,omitempty"`
	LastName       string `db:"Last_Name" json:"lastName,omitempty"`
	Email          string `db:"Email" json:"email,omitempty"`
	ProfilePicture string `db:"Profile_Picture" json:"profilePicture,omitempty"`
	Admin          bool   `db:"Admin" json:"admin,omitempty"`
	Password       string `db:"Password" json:"password,omitempty"`
}

type Comment struct {
	ID       int    `db:"ID" json:"id"`
	PostID   int    `db:"Post_ID" json:"postID"`
	AuthorID int    `db:"Author_ID" json:"authorID"`
	Content  string `db:"Content" json:"content"`
}

type Post struct {
	ID         int          `db:"ID" json:"id"`
	AuthorID   int          `db:"Author_ID" json:"authorID"`
	Title      string       `db:"Title" json:"title"`
	ImageUrl   string       `db:"Image_URL" json:"imageURL,omitempty"`
	Content    string       `db:"Content" json:"content"`
	Archived   bool         `db:"Archived" json:"archived,omitempty"`
	UploadDate sql.NullTime `db:"Upload_Date" json:"uploadDate,omitempty"`
	Slug       string       `db:"Slug" json:"slug,omitempty"`
}

type UserService interface {
	CreateUser(*User) error
	GetUsers() ([]User, error)
	GetUserByID(int) (User, error)
	GetUserByEmail(string) (User, error)
	UpdateUser(*User) error
	DeleteUser(int) error
	GetRecordCount() (int64, error)
}

type PostService interface {
	GetPosts(int, int, bool) ([]Post, error)
	GetPostById(int) (Post, error)
	GetPostBySlug(string) (Post, error)
	CreatePost(*Post) error
	UpdatePost(*Post) (int64, error)
	DeletePost(int) error
	GetRecordCount(bool) (int64, error)
}

type DBService struct {
	PostStore PostService
	UserStore UserService
}

func (p *Post) makeSlug() {
	t := strings.Split(p.Title, " ")
	slug := ""
	for _, s := range t {
		slug += s + "-"
	}
	p.Slug = fmt.Sprintf("%s%v-%v-%v", slug, time.Now().Year(), time.Now().Month(), time.Now().Day())
}
