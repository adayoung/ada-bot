package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetJSON(url string, v interface{}) error {
	if response, err := Get(url); err == nil {
		defer response.Body.Close() // So _this_ is why we had #19? I dunno :o
		if response.StatusCode == 200 {
			if data, err := ioutil.ReadAll(response.Body); err == nil {
				if err := json.Unmarshal([]byte(data), &v); err != nil {
					return err // Error at json.Unmarshal() call
				}
			} else {
				return err // Error at ioutil.ReadAll() call
			}
		} else {
			return fmt.Errorf(fmt.Sprintf("Non-OK status for: %s", url)) // <--
		}
	} else {
		return err // Error at client.Do() call
	}
	return nil
}
