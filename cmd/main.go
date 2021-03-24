package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/muazamkamal/trello-burndown/pkg/server"
	"github.com/muazamkamal/trello-burndown/pkg/trello"
	"github.com/spf13/viper"
)

func init() {
	binaryPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(binaryPath)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.WatchConfig()
	viper.SetDefault("http.readOnly", false)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	go server.Start()
	trello.Start()
}
