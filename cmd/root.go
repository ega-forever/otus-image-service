package cmd

import (
	"github.com/ega-forever/otus-image-service/internal/routes"
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
		// log.Info(port)

		r := mux.NewRouter()
		r.Use(mux.CORSMethodMiddleware(r))
		// r.Use(app.LoggingMiddleware)

		routes.SetImageRouter(r)

		err := http.ListenAndServe(":"+port, r)

		if err != nil {
			log.Panic(err)
		}

	},
}
