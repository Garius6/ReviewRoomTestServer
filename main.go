package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

var movies []Movie = []Movie{
	{0, "Movie 0", "static/example-uuid-1.jpg"},
	{1, "Movie 1", "static/example-uuid-2.jpg"},
	{2, "Movie 2", "static/example-uuid-3.jpg"},
	{3, "Movie 3", "static/example-uuid-4.jpg"},
}

var users map[string]User = map[string]User{}

var comments map[Movie][]Comment = make(map[Movie][]Comment)

var currentCommentId float64 = 0

func main() {
	getLocalIp()

	r := mux.NewRouter()
	r.HandleFunc("/auth/token", getTokenPair).Methods("GET")

	r.HandleFunc("/auth/token/refresh", logged(refreshToken)).Methods("POST")

	r.HandleFunc("/movie/{id}", logged(authorized(getMovie))).Methods("GET")

	r.HandleFunc("/movies", getMovies).Methods("GET")

	r.HandleFunc("/movie/{id}/comment", logged(authorized(createComment))).Methods("POST")

	r.HandleFunc("/movie/{id}/comments", getComments).Methods("GET")

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.Handle("/", r)
	logrus.Fatal(http.ListenAndServe(":8000", nil))
}

func logged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info(r.URL)
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Warn("Error reading request body: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		logrus.Info(string(buf))
		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader
		next(w, r)
	}
}

func authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ValidateUserToken(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")); err != nil {
			logrus.Warn(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Refresh token")
	var refreshToken string
	defer r.Body.Close()
	refreshTokenJSON, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(refreshTokenJSON, &refreshToken)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Info("validating token ", refreshToken)
	rc, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("Creating new token for ", string(refreshToken))
	newToken, err := GenerateUserToken(users[rc.Username])
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenPair := TokenPair{AccessToken: newToken, RefreshToken: refreshToken}
	tokenPairJSON, err := json.Marshal(tokenPair)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(tokenPairJSON); err != nil {
		logrus.Fatal(err)
	}
}

func getTokenPair(w http.ResponseWriter, r *http.Request) {
	var rUser User
	params := r.URL.Query()
	rUser.Username = params["username"][0]
	rUser.Password = params["password"][0]

	_, ok := users[rUser.Username]
	if !ok {
		createUser(rUser)
	}
	user := users[rUser.Username]

	logrus.Info("User = ", user)
	tokens, err := GenerateTokenPair(user)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokensJSON, err := json.Marshal(tokens)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info(fmt.Sprintf("%+v", string(tokensJSON)))
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(tokensJSON); err != nil {
		logrus.Fatal(err)
	}

}

func createUser(user User) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Warn("Hashing failed")
	}

	user.Password = string(hashedPassword)
	logrus.Info(fmt.Sprintf("Creating user %v", user))
	users[user.Username] = user
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.ParseFloat(idVar, 64)
	if err != nil {
		logrus.Warn("Atoi error ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	movie := movies[int(id)]
	movieJSON, err := json.Marshal(movie)

	logrus.Debug(movieJSON)

	if err != nil {
		logrus.Warn("Marshaling error ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(movieJSON)
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
