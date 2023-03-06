package data

import (
	"database/sql"
	"fmt"
	"reflect"
)

type MysqlUserStore struct {
	DB *sql.DB
}

// Creates a new User on the database
func (u *MysqlUserStore) CreateUser(user *User) error {
	qString := "INSERT INTO User (First_Name, Last_Name, Email, Profile_Picture, Admin, Password) VALUES (?,?,?,?,?,?);"
	stmt, err := u.DB.Prepare(qString)
	if err != nil {
		return err
	}

	user.hashPassword()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ProfilePicture, user.Admin, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Gets all users from the database
func (u *MysqlUserStore) GetUsers() ([]User, error) {
	rows, err := u.DB.Query("SELECT ID, First_Name, Last_Name, Email, Profile_Picture FROM User")
	var res []User
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		us := User{}
		err = rows.Scan(&us.ID, &us.FirstName, &us.LastName, &us.Email, &us.ProfilePicture, &us.Admin, &us.Password)
		if err != nil {
			return []User{}, err
		}
		res = append(res, us)
	}
	return res, nil
}

// Gets a user by ID
func (u *MysqlUserStore) GetUserByID(id int) (User, error) {
	qString := "SELECT * FROM User WHERE ID=?"
	stmt, err := u.DB.Prepare(qString)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return User{}, err
	}

	us := User{}

	err = stmt.QueryRow(id).Scan(&us.ID, &us.FirstName, &us.LastName, &us.Email, &us.ProfilePicture, &us.Admin, &us.Password)
	if err != nil {
		return User{}, err
	}
	return us, nil
}

// Gets a user by Email
func (u *MysqlUserStore) GetUserByEmail(email string) (User, error) {
	qString := "SELECT * FROM User WHERE Email=?"
	stmt, err := u.DB.Prepare(qString)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}

	us := User{}

	err = stmt.QueryRow(email).Scan(&us.ID, &us.FirstName, &us.LastName, &us.ProfilePicture, &us.Email, &us.Admin, &us.Password)
	if err != nil {
		return User{}, err
	}

	return us, nil
}

// Updates a user by ID
func (u *MysqlUserStore) UpdateUser(us *User) error {
	v := reflect.ValueOf(*us)
	t := v.Type()
	qString := "UPDATE User SET "
	setString := ""
	expected := make([]interface{}, 0)
	for i := 1; i < v.NumField(); i++ {
		prevState := setString
		if v.Field(i).Interface() != "" && v.Field(i).Interface() != 0 {
			setString += fmt.Sprintf(" %v=?", t.Field(i).Tag.Get("db"))
			expected = append(expected, v.Field(i).Interface())
			if i < v.NumField()-3 && setString != prevState {
				setString += ", "
			}
		}
	}

	qString += setString + " WHERE ID=?"
	expected = append(expected, us.ID)
	if setString != "" {
		stmt, err := u.DB.Prepare(qString)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(expected...)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteUser updates a user to reflect that it has been deleted, while not actually deleting the row from the database in order to prevent unexpected behaviour from
// cascading deletes or updates to Posts.
func (u *MysqlUserStore) DeleteUser(id int) error {
	us := User{ID: id, FirstName: "Deleted", LastName: "Deleted", ProfilePicture: "Deleted", Email: "Deleted", Admin: false}
	err := u.UpdateUser(&us)
	if err != nil {
		return err
	}
	return nil
}

func (u *MysqlUserStore) GetRecordCount() (int64, error) {
	var c int64
	qString := "SELECT COUNT(*) FROM User"

	err := u.DB.QueryRow(qString).Scan(&c)
	if err != nil {
		return -1, err
	}

	return c, nil
}
