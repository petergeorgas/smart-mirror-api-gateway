package main

import (
	"fmt"
	"log"
	"net/http"
	"smartmirror/handlers"

	"github.com/rs/cors"
)

func main() {

	h := handlers.NewFaceRecognitionHandler()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://smartmirror.website"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	http.Handle("/face", c.Handler(h))

	fmt.Println("Server listening on port 8080!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
