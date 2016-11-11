package util

import (
	"io/ioutil"
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration with viper.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if viper.GetBool("debug") {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}
