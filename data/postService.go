package data

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Post store acts handles data access for Post objects.
type PostStore struct {
	DB *sql.DB
}

func (p *PostStore) GetPosts(lim int) ([]Post, error) {
	if lim == 0 {
		lim = 5
	}
	var res []Post
	qString := "SELECT MAX(ID) FROM Post LIMIT " + fmt.Sprint(lim) + ";"
	rows, err := p.DB.Query(qString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ps := Post{}
		err = rows.Scan(&ps.ID, &ps.AuthorID, &ps.Content, &ps.ImageUrl, &ps.Archived)
		if err != nil {
			return nil, err
		}
		res = append(res, ps)
	}
	return res, nil
}

func (p *PostStore) GetPostById(id int) (Post, error) {
	qString := fmt.Sprintf("SELECT * FROM Post WHERE ID=%d", id)
	rows, err := p.DB.Query(qString)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()

	ps := Post{}

	rows.Next()

	err = rows.Scan(&ps.ID, &ps.AuthorID, &ps.ImageUrl, &ps.Content, &ps.Archived)
	if err != nil {
		return Post{}, err
	}

	return ps, nil
}

func (p *PostStore) CreatePost(ps *Post) (Post, error) {
	qString := fmt.Sprintf("INSERT INTO Post (Author, Image_URL, Content, Archived) VALUES (%d, '%s', '%s', %v);", ps.AuthorID, ps.ImageUrl, ps.Content, ps.Archived)
	rows, err := p.DB.Query(qString)
	if err != nil {
		return Post{}, err
	}

	res := Post{}

	rows.Next()

	err = rows.Scan(&res.ID, &res.AuthorID, &res.ImageUrl, &res.Content, &res.Archived)
	if err != nil {
		return Post{}, err
	}

	return res, nil
}

func (p *PostStore) UpdatePost(ps *Post) (int64, error) {
	v := reflect.ValueOf(*ps)
	t := v.Type()
	qString := "UPDATE Post\nSET"
	setString := ""

	for i := 1; i < v.NumField(); i++ {
		prevState := setString
		if v.Field(i).Interface() != "" && v.Field(i).Interface() != 0 {
			setString += fmt.Sprintf(" %v = '%v'", t.Field(i).Tag.Get("db"), v.Field(i).Interface())

			if i < v.NumField()-1 && setString != prevState {
				setString += ",\n"
			}
		}
	}

	qString += setString + fmt.Sprintf("\nWHERE ID = %d;", ps.ID)
	fmt.Print(qString + "\n")
	if setString != "" {
		res, err := p.DB.Exec(qString)
		if err != nil {
			return -1, err
		}
		return res.LastInsertId()
	}
	return -1, nil
}

func (p *PostStore) DeletePost(id int) {
	qString := fmt.Sprintf("UPDATE Post Archived = 'TRUE' WHERE ID=%d", id)
	_, err := p.DB.Exec(qString)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
