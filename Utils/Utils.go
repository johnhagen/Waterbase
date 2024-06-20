package Utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ReadHeader(r *http.Request) map[string]interface{} {

	//fmt.Println(r.Header)

	data := make(map[string]interface{})

	data["adminkey"] = r.Header.Get("Adminkey")
	data["auth"] = r.Header.Get("Auth")
	data["servicename"] = r.Header.Get("Servicename")
	data["collectionname"] = r.Header.Get("Collectionname")
	data["documentname"] = r.Header.Get("Documentname")

	return data
}

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

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func IsString(s interface{}) bool {
	if _, ok := s.(string); ok {
		return true
	}
	return false
}
