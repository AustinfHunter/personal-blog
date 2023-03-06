package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AustinfHunter/blog/server/data"
)

// createSuperUserFlags creates a new admin user on the database using commandline arguments.
func createSuperUserFlags(db *data.DBService) {
	suCmd := flag.NewFlagSet("create-superuser", flag.ExitOnError)
	suFirstName := suCmd.String("fname", "", "fname")
	suLastName := suCmd.String("lname", "", "lname")
	suEmail := suCmd.String("email", "", "email")
	suPassword := suCmd.String("password", "", "password")

	if len(os.Args) < 2 {
		return
	}

	fmt.Println("Attempting to create new super-user")
	if os.Args[1] == "create-superuser" {
		suCmd.Parse(os.Args[2:])
		User, err := db.UserStore.GetUserByEmail(*suEmail)
		if err != nil {
			return
		}

		if User.ID != 0 {
			return
		}

		err = db.UserStore.CreateUser(&data.User{FirstName: *suFirstName, LastName: *suLastName, Email: *suEmail, Password: *suPassword, Admin: true})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Successfully create new Super User:\nFirst Name: %s\nLast Name: %s\nEmail: %s", *suFirstName, *suLastName, *suEmail)
	}
}

// createSuperUserFlags creates a new admin user on the database using environment variables
func createSuperUserEnv(db *data.DBService) {
	fmt.Println("Attempting to create new super-user")
	fname := os.Getenv("SUFNAME")
	lname := os.Getenv("SULNAME")
	email := os.Getenv("SUEMAIL")
	pass := os.Getenv("SUPASSWORD")
	User, err := db.UserStore.GetUserByEmail(email)
	if err != nil {
		return
	}

	if User.ID != 0 {
		return
	}

	err = db.UserStore.CreateUser(&data.User{FirstName: fname, LastName: lname, Email: email, Password: pass, Admin: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully create new Super User:\nFirst Name: %s\nLast Name: %s\nEmail: %s", fname, lname, email)
}
