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
		//Check if the path is at the root
		if req.URL.Path != "/" {
			//If the path is not at the root, check if it corresponds to an existing resource
			fPath := dir + strings.TrimPrefix(path.Clean(req.URL.Path), "/")
			_, err := os.Stat(fPath)
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("err: %v\n", err)
					return
				}
				//If the path does not correspond to an existing resource, set the path to the root path before handling
				req.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
