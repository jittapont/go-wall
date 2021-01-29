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
	query := getQuery()
	log.Printf("Query : %#v", query)
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
	p, err := un.SearchPhotos(query, viper.GetInt("Page"), viper.GetInt("PerPage"), viper.GetString("Orientation"))
	if err != nil {
		log.Fatalf("Error in getting images from unsplash -> %v\n", err)
	}
	randomPhoto := getRandomPhoto(p)
	u := randomPhoto.URL.Raw
	err = wallpaper.SetFromURL(u)
	if err != nil {
		log.Fatalf("Error in setting wallpaper -> %v\n", err)
	}
}
