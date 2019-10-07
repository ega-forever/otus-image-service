package routes

import (
	"encoding/json"
	"github.com/ega-forever/otus-image-service/internal/messages"
	"github.com/gorilla/mux"
	"net/http"
)

func SetImageRouter(r *mux.Router) {
	s := r.PathPrefix("/images").Subrouter()
	s.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		imageUrl := request.URL.Query().Get("url")

		message := messages.DefaultResponse{Status: 1}
		marshaled, _ := json.Marshal(message)
		_, _ = response.Write(marshaled)
	}).Methods(http.MethodGet)

	// s.HandleFunc("/", AddProductsHandler).Methods(http.MethodPost)
	// s.HandleFunc("/", RemoveProductsHandler).Methods(http.MethodDelete)
}
