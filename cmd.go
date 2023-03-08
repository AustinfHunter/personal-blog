package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AustinfHunter/blog/server/data"
)

// createSuperUserFlags creates a new admin user on the database using environment variables
func createSuperUserEnv(db *data.DBService) {
	fmt.Println("Attempting to create new super-user")
	fname := os.Getenv("SUFNAME")
	lname := os.Getenv("SULNAME")
	email := os.Getenv("SUEMAIL")
	pass := os.Getenv("SUPASSWORD")
	User, err := db.UserStore.GetUserByEmail(email)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		fmt.Printf("Could not create new super user:: err: %v\n", err)
		return
	}

	if User.ID != 0 {
		println("Could not create new super user, that user already exists.")
		return
	}

	err = db.UserStore.CreateUser(&data.User{FirstName: fname, LastName: lname, Email: email, Password: pass, Admin: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully create new Super User:\nFirst Name: %s\nLast Name: %s\nEmail: %s", fname, lname, email)
}
