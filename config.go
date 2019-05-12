package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("AuthExpiryDays", 2)
	viper.SetDefault("AuthSecret", "NOrOmqZ2")
	viper.SetDefault("SessionSecret", "u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")
}

func init() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		// viper.AddConfigPath(filepath.Dir(dirname))
		viper.AddConfigPath(".") // Optionally look for config in the working directory.
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
		})
		if err := viper.ReadInConfig(); err != nil {
			log.Panicf("Fatal error config file: %s \n", err)
		}
	} else {
		viper.AutomaticEnv()
	}
}
