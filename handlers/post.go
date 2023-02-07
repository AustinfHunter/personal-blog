package handlers

import (
	"net/http"

	"github.com/AustinfHunter/blog/server/data"
)

func PopulatePosts(w http.ResponseWriter, req *http.Request, db *data.DBService) {
	w.Header().Set("Content-type", "application/json")
}

func getPost(w http.ResponseWriter, req *http.Request, db *data.DBService) {
	w.Header().Set("Content-type", "application/json")
}

func getComments() {

}
