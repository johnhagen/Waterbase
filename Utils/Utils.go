package Utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadFromJSON(r *http.Request) (map[string]interface{}, error) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return data, nil
}
