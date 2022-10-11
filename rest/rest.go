package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiResponse[T any] struct {
	Docs []T
}

// Get an endpoint and decode the response into a list of T's
func GetAndDecode[T any](url string) ([]T, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GetAndDecode failed to contact %v: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("GetAndDecode received error code %v from %v", response.StatusCode, url)
	}

	topLevelContainer := new(apiResponse[T])
	err = json.NewDecoder(response.Body).Decode(topLevelContainer)
	if err != nil {
		return nil, fmt.Errorf("GetAndDecode failed to decode response from %v: %w", url, err)
	}

	return topLevelContainer.Docs, nil
}

// Get an endpoint and decode the response into a single T
func GetAndDecodeSingle[T any](url string) (T, error) {
	all, err := GetAndDecode[T](url)
	if err != nil {
		return *new(T), fmt.Errorf("GetAndDecodeSingle failed: %w", err)
	}
	if len(all) != 1 {
		return *new(T), fmt.Errorf("GetAndDecodeSingle found %v reponses from %v, expected 1", len(all), url)
	}
	return all[0], nil
}

// Useful debugging tool, gets an endpoint and prints the response body
func GetAndPrint(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to contact %s: %w", url, err)
	}
	defer response.Body.Close()

	str, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read %s: %w", url, err)
	}

	fmt.Println(string(str))
	return nil
}
