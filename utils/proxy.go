package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetHtmlByProxy(pageUrl string) (string, error) {
	params := url.Values{}
	params.Add("api_key", os.Getenv("PROXY_API_KEY"))
	params.Add("url", pageUrl)
	url := "https://proxy.scrapeops.io/v1/?" + params.Encode()
	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
