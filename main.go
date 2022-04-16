package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Movie struct {
	Id        float64 `json:"id"`
	Name      string  `json:"name"`
	PosterUrl string  `json:"poster_url"`
}

func (m Movie) ToString() string {
	return fmt.Sprintf("{%f, %s,%s}", m.Id, m.Name, m.PosterUrl)
}

type Comment struct {
	Id   float64 `json:"id"`
	Text string  `json:"text"`
}

var movies []Movie = []Movie{
	{0, "Movie 0", "static/example-uuid-1.jpg"},
	{1, "Movie 1", "static/example-uuid-2.jpg"},
	{2, "Movie 2", "static/example-uuid-3.jpg"},
	{3, "Movie 3", "static/example-uuid-4.jpg"},
}

var comments map[Movie][]Comment = make(map[Movie][]Comment)

var currentId float64 = 0

func getLocalIp() {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if strings.Contains(ip.String(), "192") {
				fmt.Print("Host = ")
				fmt.Println(ip)
			}
		}
	}
}

func main() {
	getLocalIp()

	r := mux.NewRouter()
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")

	r.HandleFunc("/movies", getMovies).Methods("GET")

	r.HandleFunc("/movie/{id}/comment", createComment).Methods("POST")

	r.HandleFunc("/movie/{id}/comments", getComments).Methods("GET")

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.ParseFloat(idVar, 64)
	if err != nil {
		logrus.Fatal("Atoi error %s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	movie := movies[int(id)]
	movieJSON, err := json.Marshal(movie)

	logrus.Debug(movieJSON)

	if err != nil {
		logrus.Fatal("Marshaling error %s", err.Error())
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
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id, err := strconv.ParseFloat(vars["id"], 64)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	var newComment Comment
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newComment)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	idIdx := int(id)
	newComment.Id = currentId
	currentId++
	_, ok := comments[movies[idIdx]]
	if ok {
		comments[movies[idIdx]] = append(comments[movies[idIdx]], newComment)
	} else {
		comments[movies[idIdx]] = []Comment{newComment}
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	movieId, err := strconv.ParseFloat(vars["id"], 64)
	if err != nil {
		logrus.Fatal(err.Error())
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
		logrus.Fatal(err.Error())
	}
	fmt.Fprintf(w, string(commentsJSON))
}
