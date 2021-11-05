package main

import (
	"errors"
	"flag"
	"fmt"
	"go-wall/unsplash"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

func savePhoto(outputDir, link string) error {
	rand.Seed(time.Now().UnixNano())
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New("non-200 status code")
	}

	file := filepath.Join(outputDir, "wallpaper", fmt.Sprintf("%v.jpeg", rand.Intn(100000)))
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, res.Body)
	if err != nil {
		return err
	}

	return nil
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

	err = savePhoto("", p.URL.Raw)
	if err != nil {
		log.Fatalf("error in saving wallpaper: %s\n", err.Error())
	}

}
