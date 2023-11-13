# Ghibli Movie Fetcher Dojo

The repository is meant as a small dojo, with two levels: an advanced one where few instructions are given, and a more beginner-friendly one with a more detailed, step by step guide.
A solution is given in the `solution` folder.

The goal is to implement a small demo CLI application.
It should be written in Go and utilize the https://ghibliapi.vercel.app/ , which is a free API exposing data about the movies from the Ghibli studio.

When it starts, the application should make a call to fetch the data of all movies, then display a list with the title and release year of each.
When one movie is selected, the application should display its synopsis then exit.

### High-level instructions

After running the application in the `solution` folder to ensure your environment works fine, implement your own clone of the application, using libraries like `github.com/charmbracelet/bubbles` to handle the CLI UI part and the standard library for the rest.
Don't forget to write unit tests - `github.com/go-test/deep` and `https://github.com/stretchr/testify` might be useful.

### Step by step instructions

1. Run the application in the solution folder: `pushd solution && go run . && popd`
2. Test access to the Ghibli API we'll consume: `curl https://ghibliapi.vercel.app/films?fields=title,description,releasedate`
3. Create a new go module, with the `go mod init` command
4. Create a file named `model.go`. In it, define a type `Movie` that's a struct containing 3 fields, `Title`, `Description` and `ReleaseDate`. Add json tags so that you can parse that type from the response of the Ghibli movie API.
5. In a file named `movie_fetcher.go`, define a struct type `MovieFetcher`, with a `GetMovies` method. It should take no arguments, and return a slice of `Movies` and an `error`. Complete the method body by calling the Ghibli movie API with the help of the `http.Get` function from the standard library, and parse what's returned with `json.NewDecoder`. Don't forget to wrap errors to add more context, to close the response body (to avoid resource leak) or to check the status code of the response you get from the API.
6. Create a `main.go` file with the entrypoint of your application. In it, create a movie fetcher, use it to retrieve all movies, then print them to stdout.
7. Adapt `MovieFetcher` so that the base URL of the server isn't hardcoded anymore.
8. Write tests for the `MovieFetcher#GetMovies()` method. Use table-driven testing. You can use the following helpers:

```go
import "github.com/go-test/deep"

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
```

To mock the Ghibli API server, while still doing actual HTTP requests, use `httptest.NewServer` from the standard library.
The easiest way to do so is by letting each test case define a handlerFunc (ie, a function with the following signature: `func(w http.ResponseWriter, r *http.Request)`), and use that to initialize your `httptest.Server` with `http.HandlerFunc`.
Your test cases should also describes what outputs they're expecting, for both movie slice and error.

You should have at least test cases about:
- the Ghibli movie API returning a non 200 code
- the Ghibli movie API returning invalid JSON
- the nominal case (everything works fine)

9. Use the sample code of https://github.com/charmbracelet/bubbletea/tree/master/examples/list-simple and adapt it in a `view.go` file.
It should define a function creating a `tea.Model` which instead of displaying dishes, displays movies, and then the description of the selected one.

10. Change the code of your `main` function so that it uses what you just defined in the `view.go` file.
Optionally, fine-tune the styles of your view to ensure the description is correctly displayed.
