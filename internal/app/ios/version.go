package ios

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetVersion ...
func GetVersion() (string, error) {
	resp, err := http.Get("https://itunes.apple.com/lookup?id=944884603&country=JP")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["version"].(string), nil
}
