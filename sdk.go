package jrosenfeldLotrSdk

import (
	"fmt"

	"github.com/j-krose/jrosenfeldLotrSdk/rest"
)

var _ Sdk = (*sdk)(nil)

type Sdk interface {
	GetBooks() ([]Book, error)
	GetBook(id string) (Book, error)
}

type sdk struct {
	apiKey string
}

func NewSdk(apiKey string) Sdk {
	return sdk{apiKey}
}

// ----- /books -----

func (s sdk) GetBooks() ([]Book, error) {
	allBooks, err := rest.GetAndDecode[Book](composeUrl(book_endpoint))
	if err != nil {
		return nil, fmt.Errorf("Failed to GetBooks: %w", err)
	}
	return allBooks, nil
}

func (s sdk) GetBook(id string) (Book, error) {
	book, err := rest.GetAndDecodeSingle[Book](composeUrl(book_endpoint, id))
	if err != nil {
		return Book{}, fmt.Errorf("Failed to GetBook(%v): %w", id, err)
	}
	return book, nil
}
