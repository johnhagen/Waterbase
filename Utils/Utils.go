package Utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadFromJSON(r *http.Request) (map[string]interface{}, error) {

	data := make(map[string]interface{})

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
