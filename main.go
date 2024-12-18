package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/AustinfHunter/blog/server/cmd"
	"github.com/AustinfHunter/blog/server/data"
	"github.com/AustinfHunter/blog/server/handlers"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//Database connection setup
	mysqlU := os.Getenv("MYSQLUSER")
	mysqlPass := os.Getenv("MYSQLPASSWORD")
	mysqlHost := os.Getenv("MYSQLHOST")
	mysqlPort := os.Getenv("MYSQLPORT")
	mysqlDB := os.Getenv("MYSQLDATABASE")
	fmt.Printf("Connecting to DB at: %s:%s@tcp(%s:%s)/%s?parseTime=true\n", mysqlU, mysqlPass, mysqlHost, mysqlPort, mysqlDB)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlU, mysqlPass, mysqlHost, mysqlPort, mysqlDB))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Data access services initialization
	postStore := data.MysqlPostStore{DB: db}
	userStore := data.MysqlUserStore{DB: db}

	dbDisp := data.DBService{
		PostStore: &postStore,
		UserStore: &userStore,
	}

	//Create a new super user using environment variables
	cmd.CreateSuperUserEnv(&dbDisp)

	//Create an s3 client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	bucket := os.Getenv("AWS_BUCKET")

	//Server setup
	mux := http.NewServeMux()

	mux.Handle("/", handlers.StaticHandler(http.FileServer(http.Dir("build/")), "./build/"))

	mux.Handle("/api/posts", handlers.PopulatePosts(&dbDisp))

	mux.Handle("/api/posts/", handlers.GetPost(&dbDisp))

	mux.Handle("/api/posts/create", handlers.AddPost(&dbDisp))

	mux.Handle("/api/users/signup", handlers.SignUpHandler(&dbDisp))

	mux.Handle("/api/users/signin", handlers.SignInHandler(&dbDisp))

	mux.Handle("/api/users/authtest", handlers.AuthTestHandler(&dbDisp))

	mux.Handle("/api/admin/posts", handlers.GetAllPosts(&dbDisp))

	mux.Handle("/api/posts/update", handlers.UpdatePost(&dbDisp))

	mux.Handle("/api/upload/image", handlers.UploadImage(s3Client, bucket))

	var port string

	if os.Getenv("PORT") != "" {
		port = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	} else {
		port = ":8080"
	}

	//Start server
	t, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port %s\n", port)

	if err := http.Serve(t, mux); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
