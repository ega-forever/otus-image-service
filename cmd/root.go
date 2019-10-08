package cmd

import (
	"github.com/ega-forever/otus-image-service/internal/domain/services"
	"github.com/ega-forever/otus-image-service/internal/routes"
	"github.com/ega-forever/otus-image-service/internal/storage"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var RootCmd = &cobra.Command{
	Use:   "image service",
	Short: "service for processing images",
	Run: func(cmd *cobra.Command, args []string) {

		port := viper.GetString("REST_PORT")
		lruCache := viper.GetInt("LRU_CACHE")
		storeDir := viper.GetString("STORE_DIR")
		// log.Info(port)

		r := mux.NewRouter()
		r.Use(mux.CORSMethodMiddleware(r))
		// r.Use(app.LoggingMiddleware)

		// todo clean up dir
		st := storage.New(lruCache, storeDir)
		imageService := services.NewImageService(st)
		routes.SetImageRouter(r, imageService)

		err := http.ListenAndServe(":"+port, r)

		if err != nil {
			log.Panic(err)
		}

	},
}
