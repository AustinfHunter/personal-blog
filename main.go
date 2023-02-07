package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AustinfHunter/blog/server/data"
	"github.com/AustinfHunter/blog/server/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(172.17.0.2:3306)/blog")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	defer db.Close()

	_, err = db.Exec("CREATE TABLE if not exists User (Id INT, First_Name VARCHAR(50), Last_Name VARCHAR(50), Profile_Picture VARCHAR(150));")
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
	fmt.Println("running")

	postStore := data.PostStore{DB: db}
	userStore := data.UserStore{DB: db}

	dbDisp := data.DBService{
		PostStore: postStore,
		UserStore: userStore,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PopulatePosts(w, r, &dbDisp)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
