package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type FaceRecognitionHandler struct {
	Upgrader *websocket.Upgrader
	Conn     *websocket.Conn
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
	return &FaceRecognitionHandler{Upgrader: &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}}
}

func (h *FaceRecognitionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		conn, err := h.Upgrader.Upgrade(w, req, nil) // Upgrade to WS protocol

		if err != nil {
			fmt.Println("failed to upgrade -- client has been notified")
			fmt.Println(err)
			return
		}

		h.Conn = conn
		fmt.Println(("ws connection started!"))

	case http.MethodPost:

		var body FaceFoundRequest

		err := json.NewDecoder(req.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "invalid request body"})
			return
		}

		if h.Conn == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "WebSocket connection not found."})
			return
		}

		err = h.Conn.WriteMessage(websocket.TextMessage, []byte(body.ID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "failed to write message to WS"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FaceFoundResponse{Success: true})

	default:

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FaceFoundErrorResponse{Message: "method not allowed"})
	}

}
