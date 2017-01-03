package data

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GetLatestRelease() string {
	resp, err := http.Get("https://api.github.com/repos/meridor/perspective-backend/releases/latest")
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			var reply map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&reply)
			if err == nil {
				tagName := reply["tag_name"]
				switch tagName.(type) {
				case string:
					{
						return strings.Replace(tagName.(string), "perspective-backend-", "", -1)
					}
				}
			}
		}
	}
	return ""
}
