package routes

import (
	"context"
	"encoding/json"
	"github.com/ega-forever/otus-image-service/internal/domain/services"
	"github.com/ega-forever/otus-image-service/internal/messages"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func SetImageRouter(r *mux.Router, imageService *services.ImageService) {
	s := r.PathPrefix("/crop").Subrouter()
	s.HandleFunc("/{width}/{height}/{url:.*}", func(response http.ResponseWriter, request *http.Request) {

		vars := mux.Vars(request)
		imageUrl := vars["url"]

		width, err := strconv.Atoi(vars["width"])

		if err != nil {
			return
		}

		height, err := strconv.Atoi(vars["height"])

		if err != nil {
			return
		}

		file, headers, err := imageService.CacheToStorage(context.Background(), imageUrl, width, height)

		if err != nil {
			log.Println(err)
			message := messages.DefaultResponse{Status: 0}
			marshaled, _ := json.Marshal(message)
			response.WriteHeader(http.StatusNotFound)
			_, _ = response.Write(marshaled)
			return
		}

		for key, headerMap := range headers {
			for _, header := range headerMap {
				response.Header().Add(key, header)
			}
		}

		_, _ = response.Write(file)
	}).Methods(http.MethodGet)

}
