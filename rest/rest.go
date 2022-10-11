package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Set a high limit so that we do not need to worry about pagination.
// It would be nice to support pagination correctly, though.
// The collection with the largest number of entries is the quote
// collection with 2390 entries.
const entry_limit = 2500

// The typical format by which the API responds to our calls. See data.go for different possible `Docs`
type apiResponse[T any] struct {
	Docs  []T
	Total int
	Page  int
	Pages int
}

func appendUrlQuery(url string, key string, value string) string {
	return url + "?" + key + "=" + value
}

// For testing purposes, expose a count of how many times we have called `get`
var apiCount = 0

func ResetApiCount() {
	apiCount = 0
}

func GetApiCount() int {
	return apiCount
}

func get(url string, apiKey string) (io.ReadCloser, error) {
	url = appendUrlQuery(url, "limit", strconv.Itoa(entry_limit))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("get failed to form new request for %v: %w", url, err)
	}

	if len(apiKey) != 0 {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	apiCount += 1
	if err != nil {
		return nil, fmt.Errorf("get failed to contact %v: %w", url, err)
	}
	if response.StatusCode != 200 {
		response.Body.Close()
		return nil, fmt.Errorf("get received error code %v from %v", response.StatusCode, url)
	}
	return response.Body, nil
}

// Get an endpoint and decode the response into a list of T's
func GetAndDecode[T any](url string, apiKey string) ([]T, error) {
	body, err := get(url, apiKey)
	if err != nil {
		return nil, fmt.Errorf("GetAndDecode failed to get: %w", err)
	}
	defer body.Close()

	topLevelContainer := new(apiResponse[T])
	err = json.NewDecoder(body).Decode(topLevelContainer)
	if err != nil {
		return nil, fmt.Errorf("GetAndDecode failed to decode response from %v: %w", url, err)
	}

	// For now, report an error if our limit is set too low to accomodate the full set of entries we get back
	if topLevelContainer.Pages != 1 {
		return nil, fmt.Errorf("Too many records found for %v, %v/%v", url, topLevelContainer.Total, entry_limit)
	}

	return topLevelContainer.Docs, nil
}

// Get an endpoint and decode the response into a single T
func GetAndDecodeSingle[T any](url string, apiKey string) (T, error) {
	all, err := GetAndDecode[T](url, apiKey)
	if err != nil {
		return *new(T), fmt.Errorf("GetAndDecodeSingle failed: %w", err)
	}
	if len(all) != 1 {
		return *new(T), fmt.Errorf("GetAndDecodeSingle found %v reponses from %v, expected 1", len(all), url)
	}
	return all[0], nil
}

// Useful debugging tool, gets an endpoint and returns the response body as a string
func GetBodyAsString(url string, apiKey string) (string, error) {
	body, err := get(url, apiKey)
	if err != nil {
		return "", fmt.Errorf("GetAndPrint failed to get: %w", err)
	}
	defer body.Close()

	str, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("Failed to read %s: %w", url, err)
	}

	return string(str), nil
}
