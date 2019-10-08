package routes

import (
	"context"
	"encoding/json"
	"github.com/ega-forever/otus-image-service/internal/domain/services"
	"github.com/ega-forever/otus-image-service/internal/messages"
	"github.com/gorilla/mux"
	"net/http"
)

func SetImageRouter(r *mux.Router, imageService *services.ImageService) {
	s := r.PathPrefix("/images").Subrouter()
	s.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		imageUrl := request.URL.Query().Get("url")
		file, err := imageService.CacheToStorage(context.Background(), imageUrl)

		if err != nil {
			message := messages.DefaultResponse{Status: 0}
			marshaled, _ := json.Marshal(message)
			_, _ = response.Write(marshaled)
			return
		}

		_, _ = response.Write(file)
	}).Methods(http.MethodGet)

}
