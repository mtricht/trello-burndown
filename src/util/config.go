package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration with viper.
func InitConfig() {
	binaryPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(binaryPath)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
