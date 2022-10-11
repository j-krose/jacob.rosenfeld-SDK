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

func TestChapterFillingIsCached(t *testing.T) {
	rest.ResetApiCount()
	if len(apiKey) == 0 {
		t.Skip("Need api key for test")
	}

	sdk := NewSdk(apiKey)

	chapters, err := sdk.GetFullChapters()
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(chapters) == 0 {
		t.Errorf("Chapters was empty")
	}

	// One API call to get all the chapters, and 3 subsequent calls, each for one of the books.
	if rest.GetApiCount() != 4 {
		t.Errorf("Api made %v server calls", rest.GetApiCount())
	}
}

func TestFiltering(t *testing.T) {
	if len(apiKey) == 0 {
		t.Skip("Need api key for test")
	}

	sdk := NewSdk(apiKey)

	characters, err := sdk.GetCharacters(Matches("name", "Gandalf"))
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(characters) != 1 {
		t.Errorf("Found %v characters", len(characters))
	}
}

// I had a hard time factoring this test out to avoid repeated code, given more time I would make
// all data structs adhere to a parent type with getId and getName methods to make it easier to
// abstract things like this out
func TestEndpoints(t *testing.T) {
	if len(apiKey) == 0 {
		t.Skip("Need api key for test")
	}
	sdk := NewSdk(apiKey)

	{
		records, err := sdk.GetBooks()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("Did not find any records")
		}
		middleIndex := len(records) / 2
		middleRecord := records[middleIndex]
		middleRecordFromApi, err := sdk.GetBook(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if middleRecordFromApi != middleRecord {
			t.Errorf("%+v did not equal %+v", middleRecordFromApi, middleRecord)
		}
	}

	{
		records, err := sdk.GetMovies()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("Did not find any records")
		}
		middleIndex := len(records) / 2
		middleRecord := records[middleIndex]
		middleRecordFromApi, err := sdk.GetMovie(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if middleRecordFromApi != middleRecord {
			t.Errorf("%+v did not equal %+v", middleRecordFromApi, middleRecord)
		}
	}

	{
		records, err := sdk.GetCharacters()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("Did not find any records")
		}
		middleIndex := len(records) / 2
		middleRecord := records[middleIndex]
		middleRecordFromApi, err := sdk.GetCharacter(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if middleRecordFromApi != middleRecord {
			t.Errorf("%+v did not equal %+v", middleRecordFromApi, middleRecord)
		}
	}

	{
		records, err := sdk.GetQuotes()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("Did not find any records")
		}
		middleIndex := len(records) / 2
		middleRecord := records[middleIndex]
		middleRecordFromApi, err := sdk.GetQuote(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if middleRecordFromApi != middleRecord {
			t.Errorf("%+v did not equal %+v", middleRecordFromApi, middleRecord)
		}

		postFilled, err := sdk.FillQuote(middleRecordFromApi)
		if err != nil {
			t.Error(err)
		}
		sdkFilled, err := sdk.GetFullQuote(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if postFilled != sdkFilled {
			t.Errorf("%+v did not equal %+v", postFilled, sdkFilled)
		}
	}

	{
		records, err := sdk.GetChapters()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("Did not find any records")
		}
		middleIndex := len(records) / 2
		middleRecord := records[middleIndex]
		middleRecordFromApi, err := sdk.GetChapter(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if middleRecordFromApi != middleRecord {
			t.Errorf("%+v did not equal %+v", middleRecordFromApi, middleRecord)
		}

		postFilled, err := sdk.FillChapter(middleRecordFromApi)
		if err != nil {
			t.Error(err)
		}
		sdkFilled, err := sdk.GetFullChapter(middleRecord.Id)
		if err != nil {
			t.Error(err)
		}
		if postFilled != sdkFilled {
			t.Errorf("%+v did not equal %+v", postFilled, sdkFilled)
		}
	}
}
