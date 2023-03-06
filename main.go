package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/AustinfHunter/blog/server/data"
	"github.com/AustinfHunter/blog/server/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysqlU := os.Getenv("MYSQLUSER")
	mysqlPass := os.Getenv("MYSQLPASSWORD")
	mysqlHost := os.Getenv("MYSQLHOST")
	mysqlPort := os.Getenv("MYSQLPORT")
	mysqlDB := os.Getenv("MYSQLDATABASE")
	fmt.Printf("Connecting to DB at: %s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlU, mysqlPass, mysqlHost, mysqlPort, mysqlDB)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlU, mysqlPass, mysqlHost, mysqlPort, mysqlDB))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	postStore := data.MysqlPostStore{DB: db}
	userStore := data.MysqlUserStore{DB: db}

	dbDisp := data.DBService{
		PostStore: &postStore,
		UserStore: &userStore,
	}

	createSuperUserFlags(&dbDisp)

	mux := http.NewServeMux()

	mux.Handle("/", handlers.StaticHandler(http.FileServer(http.Dir("./static/build")), "./static/build/"))

	mux.Handle("/api/posts", handlers.PopulatePosts(&dbDisp))

	mux.Handle("/api/posts/", handlers.GetPost(&dbDisp))

	mux.Handle("/api/posts/create", handlers.AddPost(&dbDisp))

	mux.Handle("/api/users/signup", handlers.SignUpHandler(&dbDisp))

	mux.Handle("/api/users/signin", handlers.SignInHandler(&dbDisp))

	mux.Handle("/api/users/authtest", handlers.AuthTestHandler(&dbDisp))

	mux.Handle("/api/admin/posts", handlers.GetAllPosts(&dbDisp))

	mux.Handle("/api/posts/update", handlers.UpdatePost(&dbDisp))

	var port string

	if os.Getenv("PORT") != "" {
		port = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	} else {
		port = ":8080"
	}

	t, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port %s\n", port)

	if err := http.Serve(t, mux); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
