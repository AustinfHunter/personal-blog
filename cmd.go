package main

import (
	"flag"
	"os"

	"github.com/AustinfHunter/blog/server/data"
)

func createSuperUser(db *data.DBService) {
	suCmd := flag.NewFlagSet("create-superuser", flag.ExitOnError)
	suFirstName := suCmd.String("fname", "", "fname")
	suLastName := suCmd.String("lname", "", "lname")
	suEmail := suCmd.String("email", "", "email")
	suPassword := suCmd.String("pass", "", "pass")

	if os.Args[1] == "create-superuser" {
		suCmd.Parse(os.Args[2:])
		err := db.UserStore.CreateUser(&data.User{FirstName: *suFirstName, LastName: *suLastName, Email: *suEmail, Password: *suPassword})
		if err != nil {
			panic(err)
		}
	}
}
