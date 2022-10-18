package main

import (
	"fmt"
	"log"
	"net/http"
	"smartmirror/handlers"
)

func main() {

	h := handlers.NewFaceRecognitionHandler()

	http.Handle("/face", h)

	fmt.Println("Server listening on port 8080!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
