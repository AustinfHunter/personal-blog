package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/AustinfHunter/blog/server/data"
)

type postResponse struct {
	Post   data.Post  `json:"post,omitempty"`
	Author publicUser `json:"author,omitempty"`
}

type postsResponse struct {
	Posts    []data.Post `json:"posts,omitempty"`
	NumPosts int64       `json:"numPosts,omitempty"`
}

// PopulatePosts handles GET requests for posts. The number of posts returned is limited by the limit field in the request body. If there is no limit in the
// request, it will defualt to 5.
func PopulatePosts(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		w.Header().Set("Content-type", "application/json")
		var posts []data.Post
		j := json.NewEncoder(w)
		lim := 5
		offs := 0

		if req.URL.Query().Get("limit") != "" {
			qLim, err := strconv.Atoi(req.URL.Query().Get("limit"))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				lim = qLim
			}
		}

		if req.URL.Query().Get("offset") != "" {
			qOffs, err := strconv.Atoi(req.URL.Query().Get("offset"))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				offs = qOffs
			}
		}

		posts, err := db.PostStore.GetPosts(lim, offs, false)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		n, err := db.PostStore.GetRecordCount(false)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		err = j.Encode(postsResponse{posts, n})
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func GetAllPosts(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)

		if !authorizationMiddleware(req) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Content-type", "application/json")

		var posts []data.Post
		j := json.NewEncoder(w)

		posts, err := db.PostStore.GetPosts(0, 0, true)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		n, err := db.PostStore.GetRecordCount(true)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		err = j.Encode(postsResponse{posts, n})
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func GetPost(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		w.Header().Set("Content-type", "application/json")
		j := json.NewEncoder(w)

		idStr := req.URL.Query().Get("id")

		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

			p, err := db.PostStore.GetPostById(id)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

			a, err := db.UserStore.GetUserByID(p.AuthorID)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

			err = j.Encode(postResponse{p, publicUser{a.ID, a.FirstName, a.LastName}})
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			return
		}

		slug := path.Base(req.URL.Path)

		p, err := db.PostStore.GetPostBySlug(slug)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		a, err := db.UserStore.GetUserByID(p.AuthorID)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		err = j.Encode(postResponse{p, publicUser{a.ID, a.FirstName, a.LastName}})
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
	return http.HandlerFunc(fn)
}

func AddPost(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		if !authorizationMiddleware(req) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		j := json.NewEncoder(w)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}

		var post data.Post
		err = json.Unmarshal(reqBody, &post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}

		if post.Title == "" || post.Content == "" || post.AuthorID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := db.PostStore.CreatePost(&post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}

		res := struct {
			ID int64 `json:"id"`
		}{id}
		j.Encode(res)

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func UpdatePost(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		if !authorizationMiddleware(req) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}

		var post data.Post
		err = json.Unmarshal(reqBody, &post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}

		if post.Title == "" || post.Content == "" || post.AuthorID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = db.PostStore.UpdatePost(&post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("err: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
