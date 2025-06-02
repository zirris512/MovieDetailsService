package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type MovieDetails struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	Budget      int    `json:"budget"`
	Revenue     int    `json:"revenue"`
	PosterPath  string `json:"poster_path"`
	Genres      []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

func (m MovieDetails) PosterPathFull() string {
	baseUrl := "https://image.tmdb.org/t/p/w500"
	return baseUrl + m.PosterPath
}

type tmdbClient struct {
	accessToken string
	httpClient  *http.Client
}

func NewTMDbClient(accessToken string) *tmdbClient {
	return &tmdbClient{
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *tmdbClient) getMovieDetails(movieID string) (*MovieDetails, error) {
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s", movieID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Accept", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// Read the error body if available for better debugging
		errorBody, _ := io.ReadAll(response.Body)
		log.Fatalf("API returned non-OK status: %d %s. Body: %s", response.StatusCode, response.Status, errorBody)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var movieDetails MovieDetails
	err = json.Unmarshal(body, &movieDetails)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v\nJSON Body: %s", err, string(body))
	}
	movieDetails.PosterPath = movieDetails.PosterPathFull()
	return &movieDetails, nil
}
