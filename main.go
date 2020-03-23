package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySecretKey = []byte("golang")

func main() {
	router := mux.NewRouter()

	router.Path("/users").Methods(http.MethodPost).HandlerFunc(HandlerLogin)
	router.Path("/").Methods(http.MethodGet).Handler(isAuthenticated(HandlerIndex))

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email    string
		Password string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Compare hashed password
	if req.Email == "pxthang@gmail.com" && req.Password == "1234" {
		token, err := generateJWT(req.Email)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		fmt.Fprintln(w, token)
	}

}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	io.WriteString(w, `
		<h5 style="color: #111; font-family: 'Helvetica Neue', sans-serif; font-size: 100px; font-weight: bold; letter-spacing: -1px; line-height: 1; text-align: center; padding-top:50px;">JWT</h2>
		<h5 style="color: #111; font-family: 'Open Sans', sans-serif; font-size: 30px; font-weight: 300; line-height: 32px; margin: 0 0 72px; text-align: center; ">Welcome to our Project</h3>
	`)
}

// Local function
func generateJWT(username string) (string, error) {
	//
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = "Thang Pham"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		fmt.Println("Err when generating jwt")
		return "", err
	}

	return tokenString, err
}

func isAuthenticated(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid jwt token")
				}
				return mySecretKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not login yet")
		}
	})
}
