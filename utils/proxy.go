package utils

import (
	"io"
	"net/http"
	"net/url"
)

func GetHtmlByProxy(pageUrl string, render bool) (string, error) {
	params := url.Values{}
	params.Add("api_key", "d0f35644-3095-48f9-b109-9146894ae581")
	params.Add("url", pageUrl)
	params.Add("country", "us")
	if render {
		params.Add("render_js", "true")
	}
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
