package authorsearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// loadCache loads the contents of the cache file. If it doesn't
// exist, updateCache function is called.
func (website Resource) loadCache() ([]authorData, error) {
	cacheFileName := cacheFolder + "/" + website.Name + ".json"
	_, err := os.Stat(cacheFileName)
	if errors.Is(err, os.ErrNotExist) {
		err = website.updateCache()
	}
	if err != nil {
		return []authorData{}, err
	}

	stream, err := os.ReadFile(cacheFileName)
	if err != nil {
		return []authorData{}, err
	}

	data := []authorData{}
	err = json.Unmarshal(stream, &data)
	return data, err
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website Resource) updateCache() error {
	fullURL := website.BaseURL + website.QueryURL
	body, err := getResource(fullURL)
	if err != nil {
		return err
	}

	data, err := website.parseCache(body)
	if err != nil {
		return err
	}
	filteredData := filterAndDedupe(data, website.URLFilter)

	stream, err := json.MarshalIndent(filteredData, "", "    ")
	if err != nil {
		return err
	}

	cacheFileName := cacheFolder + "/" + website.Name + ".json"
	err = os.WriteFile(cacheFileName, []byte(stream), 0644)
	return err
}

func getResource(fullURL string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(fullURL)
	defer closeBody(res)

	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		return []byte{}, errors.New("Response status: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	return body, err
}

func closeBody(res *http.Response) {
	err := res.Body.Close()
	if err != nil {
		os.Stderr.WriteString("Failed to close response body: " + err.Error())
		os.Exit(1)
	}
}
