package app

import (
	"fmt"
	"github.com/paw1a/sycret-parser/internal/handler"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func Run() {
	if err := initConfig(); err != nil {
		log.Fatalf("failed to load config file")
	}

	http.HandleFunc("/api/doc", handler.DocEndpoint)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",
		viper.GetInt("port")), nil))
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
