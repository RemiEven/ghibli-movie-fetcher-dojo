package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	mf := &MovieFetcher{
		BaseURL: "https://ghibliapi.vercel.app",
	}

	movies, err := mf.GetMovies()
	if err != nil {
		slog.With("error", err).Error("failed to get movies")
		os.Exit(1)
	}

	// for _, movie := range movies {
	// 	fmt.Println(movie.Title + " (" + movie.ReleaseDate + ")")
	// }

	m := newModel(movies)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
