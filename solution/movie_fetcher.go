package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type MovieFetcher struct {
	BaseURL string
}

func (mf *MovieFetcher) GetMovies() ([]Movie, error) {
	response, err := http.Get(mf.BaseURL + "/films?fields=title,description,release_date")
	if err != nil {
		return nil, fmt.Errorf("failed to query Ghibli API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("got unexpected status " + response.Status)
	}

	var movies []Movie
	if err := json.NewDecoder(response.Body).Decode(&movies); err != nil {
		return nil, fmt.Errorf("failed to decode json returned by Ghibli API: %w", err)
	}

	return movies, err
}
