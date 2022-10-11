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

// ----- /movie
