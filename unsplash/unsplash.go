package unsplash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type GetPhotosResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Photos     []Photo `json:"results"`
}

type URL struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

type Photo struct {
	ID                     string        `json:"id"`
	CreatedAt              string        `json:"created_at"`
	UpdatedAt              string        `json:"updated_at"`
	PromotedAt             interface{}   `json:"promoted_at"`
	Width                  int           `json:"width"`
	Height                 int           `json:"height"`
	Color                  string        `json:"color"`
	BlurHash               string        `json:"blur_hash"`
	Description            interface{}   `json:"description"`
	AltDescription         string        `json:"alt_description"`
	URL                    URL           `json:"urls"`
	Likes                  int           `json:"likes"`
	LikedByUser            bool          `json:"liked_by_user"`
	CurrentUserCollections []interface{} `json:"current_user_collections"`
}

type Unsplash struct {
	BaseURL    string
	AccessKey  string
	MinTimeout time.Duration
	MaxTimeout time.Duration
	Retry      int
}

func (unsplash *Unsplash) SearchPhotos(query string, page, perPage int, orientation string) ([]Photo, error) {
	p := make([]Photo, 0)
	result := GetPhotosResponse{}
	u, err := url.Parse(unsplash.BaseURL)
	if err != nil {
		return p, nil
	}
	u.Path = path.Join(u.Path, "search", "photos")
	client, err := newClient(unsplash.MinTimeout, unsplash.MaxTimeout, unsplash.Retry)
	if err != nil {
		return p, nil
	}
	h := http.Header{
		"Authorization": []string{fmt.Sprintf("Client-ID %v", unsplash.AccessKey)},
	}
	req, err := newRequest(u, "GET", h)
	if err != nil {
		return p, nil
	}
	q := req.URL.Query()
	q.Add("query", query)
	q.Add("page", strconv.Itoa(page))
	q.Add("per_page", strconv.Itoa(perPage))
	q.Add("orientation", orientation)
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return p, nil
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return p, nil
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return p, nil
	}
	return result.Photos, nil
}

func (unsplash *Unsplash) GetRandomPhoto() (Photo, error) {
	u, err := url.Parse(unsplash.BaseURL)
	if err != nil {
		return Photo{}, err
	}

	u.Path = path.Join(u.Path, "photos", "random")

	client, err := newClient(unsplash.MinTimeout, unsplash.MaxTimeout, unsplash.Retry)
	if err != nil {
		return Photo{}, err
	}

	h := http.Header{
		"Authorization": []string{fmt.Sprintf("Client-ID %v", unsplash.AccessKey)},
	}

	req, err := newRequest(u, "GET", h)
	if err != nil {
		return Photo{}, err
	}

	q := req.URL.Query()
	q.Add("orientation", "landscape")
	q.Add("content_filter", "high")

	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return Photo{}, err
	}
	defer res.Body.Close()

	var p Photo
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return Photo{}, err
	}

	return p, nil
}
