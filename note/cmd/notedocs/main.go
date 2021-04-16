package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"noterfy/note/api/v1/transport/rest"
	_ "noterfy/note/api/v1/transport/rest/docs"
)

var port = 40001

func main() {
	r := mux.NewRouter()
	rest.RegisterDocumentationRoute(r)
	fmt.Printf("Listening to localhost:%d\n", port)
	fmt.Printf("Entry point -> http://localhost:%d/v1/index.html\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalln(err)
	}
}
