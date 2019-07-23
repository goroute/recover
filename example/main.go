package main

import (
	"log"
	"net/http"

	"github.com/goroute/recover"
	"github.com/goroute/route"
)

func main() {
	mux := route.NewServeMux()
	mux.Debug = true

	mux.Use(recover.New())

	mux.GET("/", func(c route.Context) error {
		panic("ups")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
