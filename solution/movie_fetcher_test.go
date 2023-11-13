package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestMovieFetcher(t *testing.T) {
	tests := map[string]struct {
		mockApiHandler func(http.ResponseWriter, *http.Request)
		expectedResult []Movie
		expectedErr    error
	}{
		"unexpected status": {
			mockApiHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/films" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedErr: errors.New("got unexpected status 500 Internal Server Error"),
		},
		"invalid json": {
			mockApiHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/films" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				responseBody := `[`
				if _, err := w.Write([]byte(responseBody)); err != nil {
					panic(err)
				}
			},
			expectedErr: errors.New("failed to decode json returned by Ghibli API: unexpected EOF"),
		},
		"nominal case": {
			mockApiHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/films" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				responseBody := `[{
					"title": "Nausicaä of the Valley of the Wind",
					"description": "Far in the future, after an apocalyptic conflict has devastated much of the world's ecosystem, the few surviving humans live in scattered semi-hospitable environments within what has become a 'toxic jungle.' Young Nausicaä lives in the arid Valley of the Wind and can communicate with the massive insects that populate the dangerous jungle. Under the guidance of the pensive veteran warrior, Lord Yupa, Nausicaä works to bring peace back to the ravaged planet.",
					"release_date": "1984"	
				}]`
				if _, err := w.Write([]byte(responseBody)); err != nil {
					panic(err)
				}
			},
			expectedResult: []Movie{
				{
					Title:       "Nausicaä of the Valley of the Wind",
					Description: "Far in the future, after an apocalyptic conflict has devastated much of the world's ecosystem, the few surviving humans live in scattered semi-hospitable environments within what has become a 'toxic jungle.' Young Nausicaä lives in the arid Valley of the Wind and can communicate with the massive insects that populate the dangerous jungle. Under the guidance of the pensive veteran warrior, Lord Yupa, Nausicaä works to bring peace back to the ravaged planet.",
					ReleaseDate: "1984",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(test.mockApiHandler))
			defer mockServer.Close()

			movieFetcher := &MovieFetcher{
				BaseURL: mockServer.URL,
			}

			movies, err := movieFetcher.GetMovies()

			if !ErrorEqual(err, test.expectedErr) {
				t.Errorf("unexpected error value: got [%v], wanted [%v]", err, test.expectedErr)
				return
			}

			if diff := DeepEqual(movies, test.expectedResult); diff != "" {
				t.Errorf("unexpected movies result: " + diff)
			}
		})
	}
}

// ErrorEqual checks whether two errors have the same message (or are both nil)
func ErrorEqual(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}
	if (actual != nil && expected == nil) || (actual == nil && expected != nil) {
		return false
	}
	if actual.Error() != expected.Error() {
		return false
	}
	return true
}

// DeepEqual compares a and b and returns a formatted string explaining the differences if there are any
func DeepEqual(actual, expected interface{}) string {
	diff := deep.Equal(actual, expected)
	diffHeader := "difference(s) found between actual and expected:"
	switch len(diff) {
	case 0:
		return ""
	case 1:
		return diffHeader + " " + diff[0]
	default:
		return diffHeader + "\n- " + strings.Join(diff, "\n- ")
	}
}
