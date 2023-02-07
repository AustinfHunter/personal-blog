package data

import "fmt"

type User struct {
	ID             int    `db:"ID" json:"id"`
	FirstName      string `db:"First_Name" json:"firstName"`
	LastName       string `db:"Last_Name" json:"lastName"`
	Email          string `db:"Email" json:"email"`
	ProfilePicture string `db:"Profile_Picture" json:"profilePicture"`
	Admin          bool   `db:"Admin" json:"admin"`
}

type Comment struct {
	ID       int    `db:"ID" json:"id"`
	PostID   int    `db:"Post_ID" json:"postID"`
	AuthorID int    `db:"Author_ID" json:"authorID"`
	Content  string `db:"Content" json:"content"`
}

type Post struct {
	ID       int    `db:"ID" json:"id"`
	AuthorID int    `db:"Author_ID" json:"authorID"`
	ImageUrl string `db:"Image_URL" json:"imageURL"`
	Content  string `db:"Content" json:"content"`
	Archived bool   `db:"Archived" json:"archived"`
}

type DBService struct {
	PostStore PostStore
	UserStore UserStore
}

func (u *User) ToString() string {
	return fmt.Sprintf("\n%d\n%s\n%s\n%s\n%s", u.ID, u.FirstName, u.LastName, u.Email, u.ProfilePicture)
}
