package jrosenfeldLotrSdk

const base_url = "https://the-one-api.dev/v2/"

func composeUrl(segments ...string) string {
	url := base_url
	for _, segment := range segments {
		url += segment
	}
	return url
}

// ----- /book ------
const book_endpoint = "book/"

type Book struct {
	Id   string `json:"_id"`
	Name string
}

// ----- /movie -----
const movie_endpoint = "movie/"

type Movie struct {
	Id                      string `json:"_id"`
	Name                    string
	RuntimeInMinutes        int
	BudgetInMillions        int
	AcademyAwardNominations int
	AcademyAwardWins        int
	RottenTomatoesScore     float32
}

// ----- /character -----
const character_endpoint = "character/"

type Character struct {
	Id      string `json:"_id"`
	Name    string
	Birth   string
	Death   string
	Hair    string
	Gender  string
	Height  string
	Realm   string
	Spouse  string
	Race    string
	WikiUrl string
}

// ----- /quote -----
const quote_endpoint = "quote/"

type Quote struct {
	Id          string `json:"_id"`
	Dialog      string
	MovieId     string `json:"movie"`
	CharacterId string `json:"character"`
}

// Represents a quote joined with its Movie and Character detials
type FullQuote struct {
	Id        string
	Dialog    string
	Movie     Movie
	Character Character
}

// ----- /chapter -----
const chapter_endpoint = "chapter/"

type Chapter struct {
	Id          string `json:"_id"`
	ChapterName string
	BookId      string `json:"book"`
}

// Represents a chapter joined with its Book details
type FullChapter struct {
	Id          string
	ChapterName string
	Book        Book
}
