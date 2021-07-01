package main

import (
	"log"

	"ap4y.me/cloud-mail/http"
)

func main() {
	s, err := http.NewServer("mail@ap4y.me")
	if err != nil {
		log.Fatal(err)
	}

	s.Addr = ":8000"
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
