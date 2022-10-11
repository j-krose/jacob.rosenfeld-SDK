Welcome to `jrosenfeldLotrSdk`, a software development kit targetting [The One API](https://the-one-api.dev/)

# Installation

1. [Grab an API key for The One API](https://the-one-api.dev/sign-up)
1. [Create a Go module in your local file system](https://go.dev/doc/tutorial/create-module)
1. Install this module:
   ```
   cd <your go module>
   go get github.com/j-krose/jrosenfeldLotrSdk
   ```
1. Use the module as yout see fit!
   ```
   package main

   import (
	   "fmt"

	   "github.com/j-krose/jrosenfeldLotrSdk"
   )

   func main() {
	   sdk := jrosenfeldLotrSdk.NewSdk("<your api key>")
		 books, err := sdk.GetBooks()
		 if err != nil {
			 fmt.Println(err.Error())
		 } else {
			 fmt.Printf("%+v\n", books)
		 }
   }
   ```

# Usage

The SDK provides access to The One API through the `Sdk` object, accesible by calling `NewSdk`.  The SDK exposes [The One API's routes](https://the-one-api.dev/documentation) through functions of the `Sdk` object; for example the `Sdk.GetBooks()` function represents the `/book` route.

Each route can be queried as a whole, or for a specific object by its id. This is represented in the SDK as two separate functions, a batch function and a singular function; for example `Sdk.GetBooks()` and `Sdk.GetBook(id string)`.

The routes which `Sdk` exposes can be seen in the [sdk.go](sdk.go):
```
type Sdk interface {
	GetBooks(...rest.UrlParameter) ([]Book, error)
	GetBook(string) (Book, error)
	GetMovies(...rest.UrlParameter) ([]Movie, error)
	GetMovie(string) (Movie, error)
  ...
}
```

## Filtering

The batch interfaces (`GetBooks`, `GetMovies`, etc) accept an optional variatic sequence of `UrlParameter` interface objects which represent the [The One API's filtering options](https://the-one-api.dev/documentation).  The `UrlParameter` interface is satisfied by the filtering objects in [filterOption.go](filterOption.go).  Usage exmaple:

```
gandalf, err := sdk.GetCharacters(jrosenfeldLotrSdk.Matches("name", "Gandalf"))
if err != nil {
  fmt.Println(err.Error())
} else if len(gandalf) != 1 {
  fmt.Println("Did not find exactly one Gandalf")
} else {
  fmt.Printf("%+v\n", gandalf[0])
}
```

## Data Types

When calling methods on the SDK, the data is packaged into data types that can be found in [data.go](data.go).  The data types are tightly coupled with the JSON returned from The One API.

### Object Filling

The `Quote` and `Chapter` objects make reference to other data fields using those field's ids. For example, the `Quote` object has a `MovieId` field containing the id of a movie.  For convenience, the SDK provides a mechanism to fill out the `Quote` and `Chapter` objects with the data referenced by its id.  The filled objects have their own data types, `FullQuote` and `FullChapter` which contain the fully joined data, rather than an reference id.  For example, the `FullQuote` object contains a `Movie` field with the full details of the movie, rather than just a `MovieId` string.

To an object to its full counterpart, the methods `FillQuote` and `FillChapter` are provided.

If the user desires this full data from the start, they can call the `GetFull...` function instead of the `Get...` function for that data; e.g. `GetFullChapters` instead of `GetChapters`.

## Testing

Tests are implemented in the [sdk_test.go](sdk_test.go) file.  By default, most tests are skipped because they require an API key.  To run the tests in full, fill in the `apiKey` variable at the top of that file.