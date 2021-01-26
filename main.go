package main

import (
	"flag"
	"fmt"
	"go-wall/unsplash"
	"log"

	"github.com/reujab/wallpaper"
	"github.com/spf13/viper"
)

func getQuery() string {
	query := flag.String("q", "", "Query for unsplash photos")
	flag.Parse()
	return *query
}

func main() {
	query := getQuery()
	log.Printf("Query : %#v", query)
	viper.SetConfigName("configs")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error in reading config file -> %v\n", err))
	}
	un := unsplash.Unsplash{
		BaseURL:    viper.GetString("BaseURL"),
		AccessKey:  viper.GetString("AccessKey"),
		MinTimeout: viper.GetInt("MinTimeout"),
		MaxTimeout: viper.GetInt("MaxTimeout"),
		Retry:      viper.GetInt("Retry"),
	}
	p, err := un.SearchPhotos(query, viper.GetInt("Page"), viper.GetInt("PerPage"), viper.GetString("Orientation"))
	if err != nil {
		panic(fmt.Errorf("Error in getting images from unsplash -> %v\n", err))
	}
	u := p[0].URL.Raw
	err = wallpaper.SetFromURL(u)
	if err != nil {
		panic(fmt.Errorf("Error in setting wallpaper -> %v\n", err))
	}
}
