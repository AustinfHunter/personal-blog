package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

func StaticHandler(fs http.Handler, dir string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			fs.ServeHTTP(w, req)
		}
		if req.URL.Path != "/" {
			fPath := dir + strings.TrimPrefix(path.Clean(req.URL.Path), "/")
			_, err := os.Stat(fPath)
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("err: %v\n", err)
					return
				}
				req.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
