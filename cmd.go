package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AustinfHunter/blog/server/data"
)

func createSuperUserFlags(db *data.DBService) {
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

func createSuperUserCons(db *data.DBService) {
	c, err := db.UserStore.GetRecordCount()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var v string
	if c == 0 {
		fmt.Println("Would you like to add a superuser to the server? Enter Y for yes or anything else to continue.")
		fmt.Scanln(&v)
		if v != "Y" {
			return
		}
		fmt.Println("Enter First Name:")
		var fname string
		fmt.Scanln(&fname)
		fmt.Println("Enter Last Name:")
		var lname string
		fmt.Scanln(&lname)
		fmt.Println("Enter Email:")
		var email string
		fmt.Scanln(&email)
		fmt.Println("Enter Password:")
		var pass string
		fmt.Scanln(&pass)

		err = db.UserStore.CreateUser(&data.User{FirstName: fname, LastName: lname, Email: email, Password: pass})
		if err != nil {
			panic(err)
		}
	}
}
