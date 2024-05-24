package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"waterbase/Auth"
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
	Name        string               `json:"name"`
	Owner       string               `json:"owner"`
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
	if _, exist := d.Services[r.Name]; exist {
		fmt.Println("Service with the name " + r.Name + " already exists")
		d.M.Unlock()
		return false
	}

	service := Service{}

	service.Collections = make(map[string]*Collection)
	service.Name = r.Name
	service.Owner = r.Owner

	d.Services[r.Name] = &service
	d.Services[r.Name].SaveService("./Save")

	fmt.Println("Created service: " + r.Name)
	d.M.Unlock()
	return true
}

func (d *DocumentDB) GetService(name string) *Service {
	d.M.Lock()
	if _, exist := d.Services[name]; !exist {
		fmt.Println("Could not find service " + name)
		d.M.Unlock()
		return nil
	}

	d.M.Unlock()
	return d.Services[name]
}

func (d *DocumentDB) DeleteService(name string) bool {
	d.M.Lock()
	if _, exist := d.Services[name]; exist {
		delete(d.Services, name)
		Auth.KeyDB.DeleteKey(name)
		os.RemoveAll("./Save/" + name + "/")
		os.Remove("./Save/" + name + "__")
		if _, exists := d.Services[name]; !exists {
			d.M.Unlock()
			return true
		}
		d.M.Unlock()
		return false
	}
	d.M.Unlock()
	return false
}

func (d *DocumentDB) NewLoadDB() {

	DEFAULT_SAVE_LOCATION := "./Save/"
	services, err := os.ReadDir(DEFAULT_SAVE_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, ser := range services {

		if strings.Contains(ser.Name(), "__") {

			serData, err := os.ReadFile(DEFAULT_SAVE_LOCATION + ser.Name())
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			sData := make(map[string]interface{})

			err = json.Unmarshal(serData, &sData)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			service := Service{}

			service.Name = sData["name"].(string)
			service.Owner = sData["owner"].(string)

			DocDB.CreateNewService(service)
			serv := DocDB.GetService(service.Name)

			serFolder, _, _ := strings.Cut(ser.Name(), "__")

			//fmt.Println("Search collection: " + DEFAULT_SAVE_LOCATION + serFolder)

			collections, err := os.ReadDir(DEFAULT_SAVE_LOCATION + serFolder)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, h := range collections {
				if strings.Contains(h.Name(), "__") {

					colPath := DEFAULT_SAVE_LOCATION + serFolder + "/" + h.Name()
					cFolder, _, _ := strings.Cut(h.Name(), "__")
					colFolder := DEFAULT_SAVE_LOCATION + serFolder + "/" + cFolder

					colData, err := os.ReadFile(colPath)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					cData := make(map[string]interface{})

					err = json.Unmarshal(colData, &cData)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					err = json.Unmarshal(colData, &cData)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					serv.CreateNewCollection(cData["name"].(string), cData["owner"].(string))

					newCol := serv.GetCollection(cData["name"].(string))

					//fmt.Println(colPath)

					//_ = newCol

					documents, err := os.ReadDir(colFolder)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					for _, l := range documents {

						fileData, err := os.ReadFile(colFolder + "/" + l.Name())
						if err != nil {
							fmt.Println(err.Error())
							return
						}

						file := make(map[string]interface{})

						err = json.Unmarshal(fileData, &file)
						if err != nil {
							fmt.Println(err.Error())
							return
						}

						newCol.CreateNewDocument(file["name"].(string), file["owner"].(string), file["content"])
					}
				}
			}
		}
	}
	fmt.Println("Finished loading saved services!")
}

func (d *DocumentDB) ReadDocDB() {

	data, err := os.ReadFile("ServiceDB")
	if err != nil {
		fmt.Println("DOCDB: " + err.Error())
		return
	}

	err = json.Unmarshal(data, &d.Services)
	if err != nil {
		fmt.Println("DOCDB: " + err.Error())
		return
	}
	fmt.Println("Read DocDB file!")
}

// -------------------------------------------------------------- SERVICE FUNCTIONS --------------------------------------------------------------------

func (s *Service) CreateNewCollection(name string, owner string) bool {

	DocDB.M.Lock()
	if _, exist := s.Collections[name]; exist {
		fmt.Println("Collection with the name " + name + " already exists")
		DocDB.M.Unlock()
		return false
	}

	var collection Collection

	collection.Documents = make(map[string]*Document)
	collection.Name = name
	collection.Owner = owner

	s.Collections[name] = &collection
	fmt.Println("Created collection: " + name)
	//DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true
}

func (s *Service) DeleteCollection(name string) bool {
	DocDB.M.Lock()
	if _, exist := s.Collections[name]; !exist {
		fmt.Println("Collection with the name " + name + " don't exist")
		DocDB.M.Unlock()
		return false
	}

	err := os.RemoveAll("./Save/" + s.Name + "/" + name + "/")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = os.Remove("./Save/" + s.Name + "/" + name + "__")
	if err != nil {
		fmt.Println(err.Error())
	}
	delete(s.Collections, name)
	fmt.Println("Deleted document: " + name + " from service: " + s.Name)
	DocDB.M.Unlock()
	return true
}

func (s *Service) GetCollection(name string) *Collection {
	DocDB.M.Lock()
	if _, exist := s.Collections[name]; !exist {
		fmt.Println(`Service: "` + s.Name + `" failed GET: ` + name + " - Does not exist")
		DocDB.M.Unlock()
		return nil
	}

	DocDB.M.Unlock()
	return s.Collections[name]
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

	savePath := path + "/" + s.Name

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
	if _, exist := c.Documents[name]; exist {
		fmt.Println("document with the name " + name + " already exists")
		DocDB.M.Unlock()
		return false
	}

	var document Document

	document.Name = name
	document.Owner = owner
	document.UpdatedBy = owner
	document.Content = content

	c.Documents[name] = &document
	fmt.Println("Created document: " + name)
	//DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true
}

func (c *Collection) DeleteDocument(name string) bool {

	DocDB.M.Lock()
	if _, exist := c.Documents[name]; !exist {
		fmt.Println("document with the name " + name + " don't exist")
		DocDB.M.Unlock()
		return false
	}

	delete(c.Documents, name)
	fmt.Println("Deleted document " + name + " from " + c.Name)
	//DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true
}

func (c *Collection) GetDocument(name string) *Document {
	DocDB.M.Lock()
	if _, exist := c.Documents[name]; !exist {
		fmt.Println("document with the name " + name + " don't exist")
		DocDB.M.Unlock()
		return nil
	}

	DocDB.M.Unlock()
	return c.Documents[name]
}

func (c *Collection) SaveCollection(path string) {

	col := make(map[string]interface{})

	col["name"] = c.Name
	col["owner"] = c.Owner
	col["lastUpdated"] = c.LastUpdated
	col["documents"] = []Document{}

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
	//DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true

}

func (d *Document) SaveDocument(path string) {

	data, err := json.Marshal(d)
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
