package image_test

import (
	"github.com/ega-forever/otus-image-service/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"regexp"
	"testing"
	"time"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("LOG_LEVEL", 4)
	viper.SetDefault("REST_PORT", "8080")
	viper.SetDefault("LRU_CACHE", 10)
	viper.SetDefault("STORE_DIR", "temp")

	viper.ReadInConfig()
	viper.AutomaticEnv()
}

func TestStartService(*testing.T) {
	t := time.NewTimer(time.Second * 5)

	go func() {
		cmd.RootCmd.Run(&cobra.Command{}, []string{""})
	}()
	<-t.C
}

func TestImageUpload(t *testing.T) {

	port := viper.GetString("REST_PORT")

	respResizer, err := http.Get("http://localhost:" + port + "/crop/600/800/d2908q01vomqb2.cloudfront.net/da4b9237bacccdf19c0760cab7aec4a8359010b0/2018/11/23/Picture1.png")

	if err != nil {
		log.Fatal(err)
	}

	respOriginal, err := http.Get("https://d2908q01vomqb2.cloudfront.net/da4b9237bacccdf19c0760cab7aec4a8359010b0/2018/11/23/Picture1.png")

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(respResizer.Header), len(respOriginal.Header))

	for key, value := range respOriginal.Header {
		isMatchKey, _ := regexp.MatchString(`/(server)|(Last-Modified)|(Content-Length)/i`, key)

		resizeHeaderElem := respResizer.Header.Get(key)

		assert.NotEmpty(t, resizeHeaderElem)

		if isMatchKey {
			resizerValue := respResizer.Header.Get(key)
			assert.Equal(t, value[0], resizerValue)
		}
	}
}

func TestWrongImageUpload(t *testing.T) {

	port := viper.GetString("REST_PORT")

	respResizer, err := http.Get("http://localhost:" + port + "/crop/600/800/" + "localhost/123")

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, respResizer.Status, "404 Not Found")

}

func TestWrongCropParams(t *testing.T) {

	port := viper.GetString("REST_PORT")

	respResizer, err := http.Get("http://localhost:" + port + "/crop/x/-12/" + "localhost/123")

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, respResizer.Status, "400 Bad Request")

}
