package main

import (
	"fmt"
)

type Collection struct {
	Id       float64 `json:"id"`
	Name     string  `json:"name"`
	AuthorId float64 `json:"author_id"`
	Movies   []Movie `json:"movies"`
}

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

type User struct {
	Id       float64
	Username string
	Password string
}
