package main

import (
	"flag"
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

	if os.Args[0] == "create-superuser" {
		suCmd.Parse(os.Args[1:])
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
	}
}
