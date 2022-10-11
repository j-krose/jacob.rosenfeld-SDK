package jrosenfeldLotrSdk

import (
	"testing"
)

const apiKey = ""

func TestSdk(t *testing.T) {
	sdk := NewSdk(apiKey)
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
}
