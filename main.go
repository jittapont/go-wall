package main

import (
	"flag"
	"go-wall/unsplash"
	"log"
	"math/rand"
	"time"

	"github.com/reujab/wallpaper"
	"github.com/spf13/viper"
)

func getQuery() string {
	query := flag.String("q", "", "Query for unsplash photos")
	flag.Parse()
	return *query
}

func getRandomPhoto(photos []unsplash.Photo) unsplash.Photo {
	rand.Seed(time.Now().UnixNano())
	return photos[rand.Intn(len(photos))] // #nosec
}

func main() {
	viper.SetConfigName("configs")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error in reading config file -> %v\n", err)
	}

	un := unsplash.Unsplash{
		BaseURL:    viper.GetString("BaseURL"),
		AccessKey:  viper.GetString("AccessKey"),
		MinTimeout: viper.GetDuration("MinTimeout"),
		MaxTimeout: viper.GetDuration("MaxTimeout"),
		Retry:      viper.GetInt("Retry"),
	}

	p, err := un.GetRandomPhoto()
	if err != nil {
		log.Fatalf("error in getting images from unsplash: %s\n", err.Error())
	}

	err = wallpaper.SetFromURL(p.URL.Raw)
	if err != nil {
		log.Fatalf("error in setting wallpaper: %s\n", err.Error())
	}
}
