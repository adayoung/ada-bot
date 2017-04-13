package ire

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getJSON(url string, v interface{}) error {
	client := &http.Client{}
	var data bytes.Buffer
	if request, err := http.NewRequest("GET", url, &data); err == nil {
		request.Header.Set("User-Agent", "ada-bot / https://github.com/adayoung/ada-bot")
		if response, err := client.Do(request); err == nil {
			if response.StatusCode == 200 {
				if data, err := ioutil.ReadAll(response.Body); err == nil {
					if err := json.Unmarshal([]byte(data), &v); err != nil {
						return err // Error at json.Unmarshal() call
					}
				} else {
					return err // Error at ioutil.ReadAll() call
				}
			} else {
				return nil // API call didn't return a player / Non-200 status
			}
		} else {
			return err // Error at client.Do() call
		}
	} else {
		return err // Error at http.NewRequest() call
	}
	return nil
}
