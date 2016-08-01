package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "port on which to start the server")
	flag.Parse()
	a := NewApi(NewStore())
	log.Printf("starting restfu server at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), a.Server()))
}
