package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"waterbase/Auth"
	CacheMem "waterbase/Cache"
	consts "waterbase/Data"
)

var DocDB DocumentDB

type DocumentDB struct {
	M        sync.Mutex
	Services map[string]*Service `json:"services"`
}

type Service struct {
	Name        string                 `json:"name"`
	Owner       string                 `json:"owner"`
	Collections map[string]*Collection `json:"collections"`
}

type Collection struct {
	ServiceName string               `json:"servicename"`
	Name        string               `json:"name"`
	Owner       string               `json:"owner"`
	AuthKey     string               `json:"auth"`
	LastUpdated string               `json:"lastUpdated"`
	Documents   map[string]*Document `json:"documents"`
}

type Document struct {
	UpdatedBy    string      `json:"updatedBy"`
	Name         string      `json:"name"`
	Owner        string      `json:"owner"`
	CreationDate string      `json:"creationDate"`
	LastUpdated  string      `json:"lastUpdated"`
	Content      interface{} `json:"content"`
}

// -------------------------------------------------------------- DB FUNCTIONS --------------------------------------------------------------------

func (d *DocumentDB) InitDB() {
	d.Services = make(map[string]*Service)
	err := os.MkdirAll("./Save", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (d *DocumentDB) CreateNewService(r Service) bool {

	d.M.Lock()
	_, err := os.ReadFile(consts.DEFAULT_SAVE + r.Name + "__")
	if err == nil {
		fmt.Println("Service already exists: " + r.Name)
		d.M.Unlock()
		return false
	}

	service := Service{}
	service.Collections = make(map[string]*Collection)
	service.Name = r.Name
	service.Owner = r.Owner
	service.SaveService(consts.DEFAULT_SAVE)

	fmt.Println("Created service: " + r.Name)
	d.M.Unlock()
	return true
}

func (d *DocumentDB) GetService(name string) *Service {
	d.M.Lock()
	var err error

	cachedData := CacheMem.Cache.Get("ser-" + name)
	if cachedData == nil {
		fmt.Println("Cache Miss")
		file, err := os.ReadFile(consts.DEFAULT_SAVE + name + "__")
		if err != nil {
			fmt.Println(err.Error())
			d.M.Unlock()
			return nil
		}
		CacheMem.Cache.Insert("ser-"+name, file)
		cachedData = &file
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(*cachedData, &data)
	if err != nil {
		fmt.Println(err.Error())
		d.M.Unlock()
		return nil
	}

	ser := Service{}

	ser.Name = data["name"].(string)
	ser.Owner = data["owner"].(string)
	ser.Collections = make(map[string]*Collection)

	d.M.Unlock()
	return &ser
}

func (d *DocumentDB) DeleteService(name string) bool {
	d.M.Lock()

	Auth.KeyDB.DeleteKey(name)
	CacheMem.Cache.Delete("ser-" + name)
	os.RemoveAll("./Save/" + name + "/")
	os.Remove("./Save/" + name + "__")

	d.M.Unlock()
	return true
}
