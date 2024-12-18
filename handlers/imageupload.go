package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type successResponse struct {
	ImageUrl string `json:"imageUrl"`
}

func UploadImage(s3Client *s3.Client, bucket string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CorsMiddleWare(&w, r)
		if !authorizationMiddleware(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		r.ParseMultipartForm(10 << 20) 
		f, h, err := r.FormFile("img")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	defer f.Close()
		uploader := manager.NewUploader(s3Client)
		res, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key: aws.String("upload-" + h.Filename),
			ContentType: aws.String(h.Header.Get("Content-Type")),
			Body: f,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j := json.NewEncoder(w)
		j.Encode(successResponse{ImageUrl: res.Location})
	})
}
