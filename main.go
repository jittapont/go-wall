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
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/micmonay/keybd_event"
	"github.com/reujab/wallpaper"
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

func closeVirtualDesktops() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	kb.HasSuper(true)
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_F4)

	for i := 0; i < 100; i++ {
		err = kb.Launching()
		if err != nil {
			return err
		}
	}

	return nil
}

type config struct {
	BaseURL    string        `required:"true" split_words:"true" default:"https://api.unsplash.com"`
	AccessKey  string        `required:"true" split_words:"true"`
	MinTimeout time.Duration `required:"true" split_words:"true" default:"10s"`
	MaxTimeout time.Duration `required:"true" split_words:"true" default:"30s"`
	Retry      int           `required:"true" split_words:"true" default:"3"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error in reading env file: %s\n", err.Error())
	}

	var conf config
	err = envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("error in reading config file: %s\n", err.Error())
	}

	un := unsplash.Unsplash{
		BaseURL:    conf.BaseURL,
		AccessKey:  conf.AccessKey,
		MinTimeout: conf.MinTimeout,
		MaxTimeout: conf.MaxTimeout,
		Retry:      conf.Retry,
	}

	p, err := un.GetRandomPhoto()
	if err != nil {
		log.Fatalf("error in getting images from unsplash: %s\n", err.Error())
	}

	if runtime.GOOS == "windows" {
		if err := closeVirtualDesktops(); err != nil {
			log.Printf("error in closing virtual desktops: %s\n", err)
		}
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
