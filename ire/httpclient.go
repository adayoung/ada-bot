package ire

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
				return errors.New(fmt.Sprintf("Non-OK status for: %s", url)) // <--
			}
		} else {
			return err // Error at client.Do() call
		}
	} else {
		return err // Error at http.NewRequest() call
	}
	return nil
}
