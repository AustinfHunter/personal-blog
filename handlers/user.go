package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/AustinfHunter/blog/server/data"
	"github.com/golang-jwt/jwt/v5"
)

type publicUser struct {
	ID        int
	FirstName string
	LastName  string
}

// SingInHandler is an http handler function that handles a user sign in request. Responds with a valid json web token in the authorization header as well as a redirect in the response body
// if the sign in is successful. Sends an error in the response body if the sign in is not successful.
func SignInHandler(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		w.Header().Set("Content-type", "application/json")

		var res authResponse
		j := json.NewEncoder(w)

		//read request body
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		//Unmarshal request body to a new User
		var u data.User
		err = json.Unmarshal(reqBody, &u)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		//get the user associated the email sent in the request
		uDB, err := db.UserStore.GetUserByEmail(u.Email)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		//Check for matching passwords. comparePasswords hashes the plaintext password that is passed as the first parameter.
		m := comparePasswords(&u, uDB.Password)

		//If the passwords match, create the proper response body. Note: The user from the database should be used for the creation of new JWT tokens to
		//gaurantee that the encoded to the token is correct.
		if m {
			t := getNewToken(&uDB)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				res = authResponse{
					"ERR::BAD TOKEN",
					"",
					"Something went wrong.",
				}
				err = j.Encode(res)
				if err != nil {
					fmt.Printf("err: %v\n", err)
				}
			}
			ts, err := t.SignedString([]byte(os.Getenv("SECRET_KEY")))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

			w.Header().Set("Authorization", ts)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Uid")
			w.Header().Set("Uid", fmt.Sprintf("%d", uDB.ID))
			return
		}

		//If the passwords do not match, create the proper response.
		res = authResponse{
			"ERR::BAD CREDENTIALS",
			"",
			"Sorry, the information you entered is not correct.",
		}
		err = j.Encode(res)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
	return http.HandlerFunc(fn)
}

// SingUpHandler is an http handler function that handles a user post request. A new User will be pushed to the database and the response body
// will contain a redirect and success message if the request can be successfully handled. Otherwise, the response body will contain a meaningful error message.
func SignUpHandler(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		var res authResponse
		j := json.NewEncoder(w)
		var u data.User

		w.Header().Set("Content-type", "application/json")

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		err = json.Unmarshal(reqBody, &u)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		if u.Password != req.Form.Get("repeatPassword") {
			res = authResponse{
				"",
				"",
				"Passwords don't match",
			}

			err = j.Encode(res)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			return
		}
		hp := getHashedPassword(&u)

		u.Password = hp

		err = db.UserStore.CreateUser(&u)
		if err != nil {
			res = authResponse{
				"ERR::COULD NOT CREATE NEW USER",
				"",
				"Something went wrong.",
			}
			err = j.Encode(res)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}

		res = authResponse{
			"",
			"",
			"Successfully create new user.",
		}

		err = j.Encode(res)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
	return http.HandlerFunc(fn)
}

func AuthTestHandler(db *data.DBService) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleWare(&w, req)
		if authorizationMiddleware(req) {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}
	return http.HandlerFunc(fn)
}
