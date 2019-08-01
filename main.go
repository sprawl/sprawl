package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eqlabs/sprawl/p2p"

	"github.com/julienschmidt/httprouter"
)

func greet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", params.ByName("lang"))
}

func main() {
	router := httprouter.New()
	p2p.Run()
	router.GET("/:lang", greet)
	log.Fatal(http.ListenAndServe(":8080", router))
}
