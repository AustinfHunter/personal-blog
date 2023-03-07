package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

// MysqlPostStore implements the PostService interace.
type MysqlPostStore struct {
	DB *sql.DB
}

// GetPosts queries the database and returns a slice of Posts.
func (p *MysqlPostStore) GetPosts(lim int, offs int, archived bool) ([]Post, error) {
	var res []Post
	qString := "SELECT * FROM Post WHERE Archived = ? OR Archived = 0 ORDER BY Upload_Date DESC"
	args := make([]interface{}, 0)
	args = append(args, archived)
	if lim != 0 {
		qString += " LIMIT ?"
		args = append(args, lim)
	}
	if offs != 0 {
		qString += " OFFSET ?"
		args = append(args, offs)
	}
	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		ps := Post{}
		err = rows.Scan(&ps.ID, &ps.AuthorID, &ps.Title, &ps.ImageUrl, &ps.Content, &ps.Archived, &ps.UploadDate, &ps.Slug)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		res = append(res, ps)
	}
	return res, nil
}

// GetPostByID queries the database for a post that matches the id in the parameter. If no post is found, an empty post will be returned.
func (p *MysqlPostStore) GetPostById(id int) (Post, error) {
	ps := Post{}
	qString := "SELECT * FROM Post WHERE ID=?"
	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return Post{}, err
	}

	err = stmt.QueryRow(qString).Scan(&ps.ID, &ps.AuthorID, &ps.Title, &ps.ImageUrl, &ps.Content, &ps.Archived, &ps.UploadDate, &ps.Slug)
	if err != nil {
		return Post{}, err
	}
	return ps, nil
}

// GetPostBySlug queries the database for a post that matches the slug passed in the parameter. If no post is found, an empty post will be returned.
func (p *MysqlPostStore) GetPostBySlug(slug string) (Post, error) {
	qString := "SELECT * FROM Post WHERE Slug=?"
	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return Post{}, err
	}
	defer stmt.Close()

	ps := Post{}

	err = stmt.QueryRow(slug).Scan(&ps.ID, &ps.AuthorID, &ps.Title, &ps.ImageUrl, &ps.Content, &ps.Archived, &ps.UploadDate, &ps.Slug)
	if err != nil {
		return Post{}, err
	}

	return ps, nil
}

// Create post adds a new post to the database and returns the id and slug if successful.
func (p *MysqlPostStore) CreatePost(ps *Post) (int64, error) {
	ps.makeSlug()
	qString := "INSERT INTO Post (Author_ID, Title, Image_URL, Content, Archived, Upload_Date, Slug) VALUES (?,?,?,?,?,?,?);"
	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(ps.AuthorID, ps.Title, "", ps.Content, ps.Archived, time.Now().Format("2006-01-02 15:04:05"), ps.Slug)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

// UpdatePost updates a Post
func (p *MysqlPostStore) UpdatePost(ps *Post) (int64, error) {
	v := reflect.ValueOf(*ps)
	t := v.Type()
	qString := "UPDATE Post SET"
	setString := ""
	expected := make([]interface{}, 0)
	for i := 1; i < v.NumField(); i++ {
		prevState := setString
		if v.Field(i).Interface() != "" && v.Field(i).Interface() != 0 {
			if t.Field(i).Tag.Get("db") == "Upload_Date" {
				continue
			}
			setString += fmt.Sprintf(" %v=?", t.Field(i).Tag.Get("db"))
			expected = append(expected, v.Field(i).Interface())
			if i < v.NumField()-3 && setString != prevState {
				setString += ", "
			}
		}
	}

	qString += setString + " WHERE ID=?"
	expected = append(expected, ps.ID)
	if setString != "" {
		stmt, err := p.DB.Prepare(qString)
		if err != nil {
			return -1, err
		}

		res, err := stmt.Exec(expected...)
		if err != nil {
			return -1, err
		}
		return res.LastInsertId()
	}
	return -1, nil
}

func (p *MysqlPostStore) DeletePost(id int) error {
	qString := "DELETE FROM Post WHERE ID=?"
	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return err
	}
	stmt.Exec(id)
	return nil
}

func (p *MysqlPostStore) GetRecordCount(archived bool) (int64, error) {
	var c int64
	qString := "SELECT COUNT(*) FROM Post WHERE Archived = 0 OR Archived = ?"

	stmt, err := p.DB.Prepare(qString)
	if err != nil {
		return -1, err
	}

	err = stmt.QueryRow(archived).Scan(&c)
	if err != nil {
		return -1, err
	}
	return c, nil
}
