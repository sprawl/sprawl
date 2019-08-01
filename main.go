package main

import {
	"fmt"
	"github.com/julienschmidt/httprouter"
}

func greet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", params.ByName("lang"))
}

func main() {
	router := httprouter.New()
	router.GET("/:lang", greet)
	router.POST("/:lang", languages.Run)
	log.Fatal(http.ListenAndServe(":8080", router))
}
