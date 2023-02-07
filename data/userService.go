package data

import (
	"database/sql"
	"fmt"
	"reflect"
)

type UserStore struct {
	DB *sql.DB
}

//Creates a new User on the database
func (u *UserStore) CreateUser(user *User) {
	qString := fmt.Sprintf("INSERT INTO User (First_Name, Last_Name, Profile_Picture) VALUES ('%s', '%s', '%s');", user.FirstName, user.LastName, user.ProfilePicture)
	_, err := u.DB.Query(qString)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
}

//Gets all users from the database
func (u *UserStore) GetUsers() []User {
	rows, err := u.DB.Query("SELECT ID, First_Name, Last_Name, Email, Profile_Picture FROM User")
	var res []User
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		us := User{}
		err = rows.Scan(&us.ID, &us.FirstName, &us.LastName, &us.Email, &us.ProfilePicture)
		if err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
		}
		res = append(res, us)
	}
	return res
}

//Gets a user by ID
func (u *UserStore) GetUserByID(id int) User {
	qString := fmt.Sprintf("SELECT * FROM User WHERE ID=%d", id)
	rows, err := u.DB.Query(qString)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}

	us := User{}

	rows.Next()

	err = rows.Scan(&us.ID, &us.FirstName, &us.LastName, &us.ProfilePicture)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}

	return us
}

//Updates a user by ID
func (u *UserStore) UpdateUser(user *User) {
	v := reflect.ValueOf(*user)
	t := v.Type()
	qString := "UPDATE User\nSET"
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

	qString += setString + fmt.Sprintf("\nWHERE ID = %d;", user.ID)
	fmt.Print(qString + "\n")
	if setString != "" {
		_, err := u.DB.Exec(qString)
		if err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
		}
	}
}

//Deletes a user by ID.
func (u *UserStore) DeleteUser(id int) {
	us := User{ID: id, FirstName: "Deleted", LastName: "Deleted", ProfilePicture: "Deleted", Email: "Deleted", Admin: false}
	u.UpdateUser(&us)
}
