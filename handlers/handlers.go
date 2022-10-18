package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type FaceRecognitionHandler struct {
	Upgrader      *websocket.Upgrader
	OutputChannel chan string
}

type FaceFoundRequest struct {
	ID string `json:"id"`
}

type FaceFoundResponse struct {
	Success bool `json:"success"`
}

type FaceFoundErrorResponse struct {
	Message string `json:"errorMessage"`
}

func NewFaceRecognitionHandler() *FaceRecognitionHandler {
	return &FaceRecognitionHandler{Upgrader: &websocket.Upgrader{}, OutputChannel: make(chan string, 2)}
}

func (h *FaceRecognitionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		conn, err := h.Upgrader.Upgrade(w, req, nil) // Upgrade to WS protocol

		if err != nil {
			fmt.Println("failed to upgrade -- client has been notified")
			return
		}

		defer conn.Close()

		for {
			id := <-h.OutputChannel

			err := conn.WriteMessage(websocket.TextMessage, []byte(id))

			if err != nil {
				fmt.Printf("error writing id %s", id)
				break
			}
		}

	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var body FaceFoundRequest

		err := json.NewDecoder(req.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "invalid request body"})
		}

		h.OutputChannel <- body.ID

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FaceFoundResponse{Success: true})

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "method not allowed"})
	}

}
