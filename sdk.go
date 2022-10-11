package jrosenfeldLotrSdk

import (
	"fmt"

	"github.com/j-krose/jrosenfeldLotrSdk/rest"
)

var _ Sdk = (*sdk)(nil)

type Sdk interface {
	GetBooks() ([]Book, error)
	GetBook(string) (Book, error)
	GetMovies() ([]Movie, error)
	GetMovie(string) (Movie, error)
	GetCharacters() ([]Character, error)
	GetCharacter(string) (Character, error)
	GetQuotes() ([]Quote, error)
	GetQuote(string) (Quote, error)
	FillQuote(Quote) (FullQuote, error)
	GetFullQuotes() ([]FullQuote, error)
	GetFullQuote(string) (FullQuote, error)
	GetChapters() ([]Chapter, error)
	GetChapter(string) (Chapter, error)
	FillChapter(Chapter) (FullChapter, error)
	GetFullChapters() ([]FullChapter, error)
	GetFullChapter(string) (FullChapter, error)
}

type sdk struct {
	apiKey string
}

func NewSdk(apiKey string) sdk {
	return sdk{apiKey}
}

// ----- /books -----

func (s sdk) GetBooks() ([]Book, error) {
	return rest.GetAndDecode[Book](composeUrl(book_endpoint), s.apiKey)
}

func (s sdk) GetBook(id string) (Book, error) {
	return rest.GetAndDecodeSingle[Book](composeUrl(book_endpoint, id), s.apiKey)
}

// ----- /movies -----

func (s sdk) GetMovies() ([]Movie, error) {
	return rest.GetAndDecode[Movie](composeUrl(movie_endpoint), s.apiKey)
}

func (s sdk) GetMovie(id string) (Movie, error) {
	return rest.GetAndDecodeSingle[Movie](composeUrl(movie_endpoint, id), s.apiKey)
}

// ----- /character -----

func (s sdk) GetCharacters() ([]Character, error) {
	return rest.GetAndDecode[Character](composeUrl(character_endpoint), s.apiKey)
}

func (s sdk) GetCharacter(id string) (Character, error) {
	return rest.GetAndDecodeSingle[Character](composeUrl(character_endpoint, id), s.apiKey)
}

// ----- /quote -----

func (s sdk) GetQuotes() ([]Quote, error) {
	return rest.GetAndDecode[Quote](composeUrl(quote_endpoint), s.apiKey)
}

func (s sdk) GetQuote(id string) (Quote, error) {
	return rest.GetAndDecodeSingle[Quote](composeUrl(quote_endpoint, id), s.apiKey)
}

func (s sdk) FillQuote(quote Quote) (FullQuote, error) {
	return s.fillQuote(quote, nil, nil)
}

// Fill a list of quotes with their character and movie details, with caching to avoid redundant server calls
func (s sdk) fillQuote(quote Quote, movies map[string]Movie, characters map[string]Character) (FullQuote, error) {
	fullQuote := FullQuote{
		Id:     quote.Id,
		Dialog: quote.Dialog,
	}

	movieCall := func(movieId string) (Movie, error) { return s.GetMovie(movieId) }
	movie, err := getDetails(quote.MovieId, movies, movieCall)
	if err != nil {
		return FullQuote{}, fmt.Errorf("fillQuote could not find movie %v: %w", quote.MovieId, err)
	}
	fullQuote.Movie = movie

	characterCall := func(characterId string) (Character, error) { return s.GetCharacter(characterId) }
	character, err := getDetails(quote.CharacterId, characters, characterCall)
	if err != nil {
		return FullQuote{}, fmt.Errorf("fillQuote could not find character %v: %w", quote.CharacterId, err)
	}
	fullQuote.Character = character

	return fullQuote, nil
}

func (s sdk) GetFullQuotes() ([]FullQuote, error) {
	quotes, err := s.GetQuotes()
	if err != nil {
		return nil, err
	}

	fullQuotes := make([]FullQuote, len(quotes))
	movies := make(map[string]Movie)
	characters := make(map[string]Character)
	for index, quote := range quotes {
		fullQuote, err := s.fillQuote(quote, movies, characters)
		if err != nil {
			return nil, err
		}
		fullQuotes[index] = fullQuote
	}

	return fullQuotes, nil
}

func (s sdk) GetFullQuote(id string) (FullQuote, error) {
	quote, err := s.GetQuote(id)
	if err != nil {
		return FullQuote{}, err
	}
	return s.FillQuote(quote)
}

// ----- /chapter -----
func (s sdk) GetChapters() ([]Chapter, error) {
	return rest.GetAndDecode[Chapter](composeUrl(chapter_endpoint), s.apiKey)
}

func (s sdk) GetChapter(id string) (Chapter, error) {
	return rest.GetAndDecodeSingle[Chapter](composeUrl(chapter_endpoint, id), s.apiKey)
}

func (s sdk) FillChapter(chapter Chapter) (FullChapter, error) {
	return s.fillChapter(chapter, nil)
}

// Fill a list of chapter with its book details, with caching to avoid redundant server calls
func (s sdk) fillChapter(chapter Chapter, books map[string]Book) (FullChapter, error) {
	fullChapter := FullChapter{
		Id:          chapter.Id,
		ChapterName: chapter.ChapterName,
	}

	bookCall := func(bookId string) (Book, error) { return s.GetBook(bookId) }
	book, err := getDetails(chapter.BookId, books, bookCall)
	if err != nil {
		return FullChapter{}, fmt.Errorf("fillChapter could not find book %v: %w", chapter.BookId, err)
	}
	fullChapter.Book = book

	return fullChapter, nil
}

func (s sdk) GetFullChapters() ([]FullChapter, error) {
	chapters, err := s.GetChapters()
	if err != nil {
		return nil, err
	}

	fullChapters := make([]FullChapter, len(chapters))
	books := make(map[string]Book)
	for index, chapter := range chapters {
		fullChapter, err := s.fillChapter(chapter, books)
		if err != nil {
			return nil, err
		}
		fullChapters[index] = fullChapter
	}

	return fullChapters, nil
}

func (s sdk) GetFullChapter(id string) (FullChapter, error) {
	quote, err := s.GetChapter(id)
	if err != nil {
		return FullChapter{}, err
	}
	return s.FillChapter(quote)
}

// ----- UTILS -----

// Pull details from either a cache or a lamda (which represents an api call)
func getDetails[T any](detailId string, detailCache map[string]T, detailFunc func(string) (T, error)) (T, error) {
	if detailCache != nil {
		detail, ok := detailCache[detailId]
		if ok {
			return detail, nil
		}
	}

	detail, err := detailFunc(detailId)
	if err != nil {
		return *new(T), fmt.Errorf("getDetails could not find detail %v: %w", detailId, err)
	}
	if detailCache != nil {
		detailCache[detailId] = detail
	}
	return detail, nil
}
