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

var movies []Movie = []Movie{
	{0, "Movie 1", "/static/example-uuid-1.jpg"},
	{1, "Movie 2", "/static/example-uuid-2.jpg"},
	{2, "Movie 3", "/static/example-uuid-3.jpg"},
	{3, "Movie 4", "/static/example-uuid-4.jpg"},
}

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

	logrus.Debug(movieJSON)

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
