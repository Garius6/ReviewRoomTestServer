package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       float64 `json:"id"`
	Name     string  `json:"name"`
	PosterId string  `json:"poster_id"`
}

func (m Movie) ToString() string {
	return fmt.Sprintf("{%f, %s,%s}", m.Id, m.Name, m.PosterId)
}

var movies []Movie = []Movie{
	{0, "Movie 1", "/static/example-uuid-1.jpg"},
	{1, "Movie 2", "/static/example-uuid-2.jpg"},
	{2, "Movie 3", "/static/example-uuid-3.jpg"},
	{3, "Movie 4", "/static/example-uuid-4.jpg"},
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movie/{id}", getMovie)

	r.HandleFunc("/movies", getMovies)

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.Atoi(idVar)
	if err != nil {
		fmt.Errorf("Atoi error %s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	movie := movies[id]
	movieJSON, err := json.Marshal(movie)

	if err != nil {
		fmt.Errorf("Marshaling error %s", err.Error())
	}
	fmt.Fprintf(w, string(movieJSON))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	moviesJSON, _ := json.Marshal(movies)
	fmt.Fprint(w, string(moviesJSON))
}