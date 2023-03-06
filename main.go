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
	mysqlPrt := os.Getenv("MYSQLPORT")
	mysqlPswd := os.Getenv("MYSQLPASSWORD")
	mysqlHst := os.Getenv("MYSQLHOST")
	mysqlDB := os.Getenv("MYSQLDATABASE")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlU, mysqlPswd, mysqlHst, mysqlPrt, mysqlDB))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	fmt.Printf("Successfully Connectect to database\n")
	defer db.Close()

	postStore := data.MysqlPostStore{DB: db}
	userStore := data.MysqlUserStore{DB: db}

	dbDisp := data.DBService{
		PostStore: &postStore,
		UserStore: &userStore,
	}

	createSuperUserCons(&dbDisp)

	mux := http.NewServeMux()

	//fs := http.FileServer(http.Dir("./static/build"))

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
		port = os.Getenv("PORT")
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
