package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"waterbase/Auth"
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
	//cLock       sync.Mutex
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
	file, err := os.ReadFile(consts.DEFAULT_SAVE + name + "__")
	if err != nil {
		fmt.Println(err.Error())
		d.M.Unlock()
		return nil
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(file, &data)
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
	os.RemoveAll("./Save/" + name + "/")
	os.Remove("./Save/" + name + "__")

	d.M.Unlock()
	return true
}

// -------------------------------------------------------------- SERVICE FUNCTIONS --------------------------------------------------------------------

func (s *Service) CreateNewCollection(name string, owner string) bool {
	fmt.Println("kekek")
	fmt.Println(s)
	fmt.Println(name)
	DocDB.M.Lock()

	_, err := os.Stat(consts.DEFAULT_SAVE + s.Name + "/" + name + "__/")
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
	fmt.Println("Created collection: " + name)
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
	delete(s.Collections, name)
	fmt.Println("Deleted document: " + name + " from service: " + s.Name)
	DocDB.M.Unlock()
	return true
}

func (s *Service) GetCollection(name string) *Collection {
	DocDB.M.Lock()
	file, err := os.ReadFile(consts.DEFAULT_SAVE + s.Name + "/" + name + "__")
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return nil
	}

	data := make(map[string]interface{})

	collection := Collection{}

	err = json.Unmarshal(file, &data)
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

// -------------------------------------------------------------- COLLECTION FUNCTIONS --------------------------------------------------------------------

func (c *Collection) CreateNewDocument(name string, owner string, content interface{}) bool {

	DocDB.M.Lock()
	_, err := os.Stat(consts.DEFAULT_SAVE + c.ServiceName + "/" + c.Name + "/" + name)
	if err == nil {
		fmt.Println("Document already exists")
		DocDB.M.Unlock()
		return false
	}

	var document Document
	document.Name = name
	document.Owner = owner
	document.UpdatedBy = owner
	document.Content = content

	document.SaveDocument(consts.DEFAULT_SAVE + c.ServiceName + "/" + c.Name)
	fmt.Println("Created document: " + name)
	DocDB.M.Unlock()
	return true
}

func (c *Collection) GetDocument(name string) *Document {
	DocDB.M.Lock()
	file, err := os.ReadFile(consts.DEFAULT_SAVE + c.ServiceName + "/" + c.Name + "/" + name)
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return nil
	}

	document := Document{}

	err = json.Unmarshal(file, &document)
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return nil
	}

	DocDB.M.Unlock()
	return &document
}

func (c *Collection) DeleteDocument(name string) bool {

	DocDB.M.Lock()
	err := os.Remove(consts.DEFAULT_SAVE + c.ServiceName + "/" + c.Name + "/" + name)
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return false
	}
	fmt.Println("Deleted document " + name + " from " + c.Name)

	DocDB.M.Unlock()
	return true
}

func (c *Collection) SaveCollection(path string) {

	col := make(map[string]interface{})

	col["name"] = c.Name
	col["owner"] = c.Owner
	col["lastUpdated"] = c.LastUpdated
	col["servicename"] = c.ServiceName
	col["auth"] = c.AuthKey
	col["documents"] = &[]Document{}

	data, err := json.Marshal(col)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path + "/" + c.Name

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

// -------------------------------------------------------------- DOCUMENT FUNCTIONS --------------------------------------------------------------------

func (d *Document) GetContent() interface{} {
	return d.Content
}

func (d *Document) SetContent(name string, content interface{}) bool {

	DocDB.M.Lock()
	d.UpdatedBy = name
	d.Content = content
	DocDB.M.Unlock()
	return true

}

func (d *Document) SaveDocument(path string) {

	temp := make(map[string]interface{})

	temp["updatedBy"] = d.UpdatedBy
	temp["name"] = d.Name
	temp["owner"] = d.Owner
	temp["creationDate"] = d.CreationDate
	temp["lastUpdated"] = d.LastUpdated
	temp["content"] = d.Content

	data, err := json.Marshal(temp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path + "/" + d.Name

	file, err := os.Create(savePath)
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
}
