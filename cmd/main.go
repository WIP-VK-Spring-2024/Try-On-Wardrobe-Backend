package main

import (
	"flag"
	"fmt"
	"log"

	"try-on/internal/pkg/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var configPath *string = flag.String("c", "config/config.json", "Specify config path")

func main() {
	viper.SetConfigFile(*configPath)

	cfg := config.Config{}

	viper.OnConfigChange(func(in fsnotify.Event) {
		tmp := config.Config{}
		err := viper.Unmarshal(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		cfg = tmp
		fmt.Printf("%+v\n", cfg)
	})

	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%+v\n", cfg)

	app := NewApp(&cfg)

	log.Fatal(app.Run())
}
