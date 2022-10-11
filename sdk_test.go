package jrosenfeldLotrSdk

import (
	"testing"

	"github.com/j-krose/jrosenfeldLotrSdk/rest"
)

// Ideally, there would be some mechanism to pick this up from a file that not version controlled
const apiKey = ""

func TestBooks(t *testing.T) {
	rest.ResetApiCount()
	sdk := NewSdk("" /* intentionally left blank to test that books can be accessed un-authenticated */)
	books, err := sdk.GetBooks()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(books) != 3 {
		t.Errorf("Got %v books", len(books))
	}

	bookId := books[0].Id
	book, err := sdk.GetBook(bookId)
	if err != nil {
		t.Errorf("%v", err)
	}
	if book.Name != books[0].Name {
		t.Errorf("Expected %v, got %v", book.Name, books[0].Name)
	}

	if rest.GetApiCount() != 2 {
		t.Errorf("Api made %v server calls", rest.GetApiCount())
	}
}

func TestFillingIsCached(t *testing.T) {
	rest.ResetApiCount()
	if len(apiKey) == 0 {
		t.Skip("Need api key for test")
	}

	sdk := NewSdk(apiKey)

	_, err := sdk.GetFullChapters()
	if err != nil {
		t.Errorf("%v", err)
	}

	// One API call to get all the chapters, and 3 subsequent calls, each for one of the books.
	if rest.GetApiCount() != 4 {
		t.Errorf("Api made %v server calls", rest.GetApiCount())
	}
}
