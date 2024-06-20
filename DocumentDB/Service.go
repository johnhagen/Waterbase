package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	CacheMem "waterbase/Cache"
	consts "waterbase/Data"
)

// -------------------------------------------------------------- SERVICE FUNCTIONS --------------------------------------------------------------------

func (s *Service) CreateNewCollection(name string, owner string) bool {
	DocDB.M.Lock()

	_, err := os.Stat(consts.DEFAULT_SAVE + s.Name + "/" + name)
	if err == nil {
		fmt.Println("Collection named: " + name + " already exist")
		DocDB.M.Unlock()
		return false
	}

	var collection Collection

	collection.Documents = make(map[string]*Document)
	collection.ServiceName = s.Name
	collection.Name = name
	collection.Owner = owner
	collection.LastUpdated = "temp"

	collection.SaveCollection(consts.DEFAULT_SAVE + s.Name + "/")
	fmt.Println("Service: " + s.Name + " - Created collection: " + name)
	DocDB.M.Unlock()
	return true
}

func (s *Service) DeleteCollection(name string) bool {
	DocDB.M.Lock()

	err := os.RemoveAll("./Save/" + s.Name + "/" + name + "/")
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return false
	}
	err = os.Remove("./Save/" + s.Name + "/" + name + "__")
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return false
	}
	CacheMem.Cache.Delete("col-" + name)
	delete(s.Collections, name)
	fmt.Println("Deleted document: " + name + " from service: " + s.Name)
	DocDB.M.Unlock()
	return true
}

func (s *Service) GetCollection(name string) *Collection {
	DocDB.M.Lock()

	cachedData := CacheMem.Cache.Get("col-" + name)
	if cachedData == nil {
		file, err := os.ReadFile(consts.DEFAULT_SAVE + s.Name + "/" + name + "__")
		if err != nil {
			fmt.Println(err.Error())
			DocDB.M.Unlock()
			return nil
		}
		cachedData = &file
		CacheMem.Cache.Insert("col-"+name, file)
	}

	data := make(map[string]interface{})

	collection := Collection{}

	err := json.Unmarshal(*cachedData, &data)
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return nil
	}

	collection.ServiceName = data["servicename"].(string)
	collection.Name = data["name"].(string)
	collection.Owner = data["owner"].(string)
	collection.LastUpdated = data["lastUpdated"].(string)
	collection.Documents = make(map[string]*Document)

	DocDB.M.Unlock()
	return &collection
}

func (s *Service) SaveService(path string) {

	ser := make(map[string]interface{})

	ser["name"] = s.Name
	ser["owner"] = s.Owner
	ser["collections"] = []Collection{}

	data, err := json.Marshal(ser)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path + s.Name

	file, err := os.Create(savePath + "__")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	file.Close()

	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
