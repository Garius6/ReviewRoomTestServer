package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrus.SetReportCaller(true)
}

var movies []Movie = []Movie{
	{0, "Movie 0", "static/example-uuid-1.jpg"},
	{1, "Movie 1", "static/example-uuid-2.jpg"},
	{2, "Movie 2", "static/example-uuid-3.jpg"},
	{3, "Movie 3", "static/example-uuid-4.jpg"},
}

var users map[string]User = map[string]User{
	"User1": {"User1", "User1Password"},
}

var comments map[Movie][]Comment = make(map[Movie][]Comment)

var currentCommentId float64 = 0

func main() {
	getLocalIp()

	r := mux.NewRouter()
	r.HandleFunc("/user/create", createUser).Methods("POST")

	r.HandleFunc("/user/login", loginUser).Methods("POST")

	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")

	r.HandleFunc("/movies", getMovies).Methods("GET")

	r.HandleFunc("/movie/{id}/comment", authorized(createComment)).Methods("POST")

	r.HandleFunc("/movie/{id}/comments", getComments).Methods("GET")

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateToken(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func returnError(w http.ResponseWriter, errorCode int) {
	w.WriteHeader(errorCode)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	savedUser, ok := users[user.Username]
	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password)); !ok || err != nil {
		logrus.Warn("User does not exist")
		returnError(w, http.StatusNotFound)
		return
	}

	token, err := generateToken(user.Username)
	if err != nil {
		logrus.Warn(err)
	}
	fmt.Fprintf(w, token)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var rUser User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&rUser)

	user, ok := users[rUser.Username]
	if ok {
		logrus.Warn(fmt.Sprintf("User %v tried to create profile again", rUser))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Warn("Hashing failed")
	}

	user.Password = string(hashedPassword)
	users[rUser.Username] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.ParseFloat(idVar, 64)
	if err != nil {
		logrus.Warn("Atoi error %s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	movie := movies[int(id)]
	movieJSON, err := json.Marshal(movie)

	logrus.Debug(movieJSON)

	if err != nil {
		logrus.Warn("Marshaling error %s", err.Error())
	}
	fmt.Fprintf(w, string(movieJSON))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	moviesJSON, _ := json.Marshal(movies)
	fmt.Fprint(w, string(moviesJSON))
}

func createComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseFloat(vars["id"], 64)
	if err != nil {
		logrus.Warn(err.Error())
	}

	var newComment Comment
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newComment)
	if err != nil {
		logrus.Warn(err.Error())
	}

	idIdx := int(id)
	newComment.Id = currentCommentId
	currentCommentId++
	_, ok := comments[movies[idIdx]]
	if ok {
		comments[movies[idIdx]] = append(comments[movies[idIdx]], newComment)
	} else {
		comments[movies[idIdx]] = []Comment{newComment}
	}
	w.WriteHeader(http.StatusOK)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	movieId, err := strconv.ParseFloat(vars["id"], 64)
	if err != nil {
		logrus.Warn(err.Error())
	}
	movieIdIdx := int(movieId)
	movieComments, ok := comments[movies[movieIdIdx]]
	if !ok {
		emptyList, _ := json.Marshal([]Comment{})
		fmt.Fprintf(w, string(emptyList))
		logrus.Info(fmt.Sprintf("Movie with id %d comments section is empty", movieIdIdx))
		return
	}
	commentsJSON, err := json.Marshal(movieComments)
	if err != nil {
		logrus.Warn(err.Error())
	}
	fmt.Fprintf(w, string(commentsJSON))
}
