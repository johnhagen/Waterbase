package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var DocDB DocumentDB

type DocumentDB struct {
	M        sync.Mutex
	Services map[string]*Service `json:"services"`
}

type Service struct {
	//sLock       sync.Mutex
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

	fmt.Println("Created service: " + r.Name)
	d.SaveDocDB()
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

func (d *DocumentDB) SaveDocDB() {

	data, err := json.Marshal(d.Services)
	if err != nil {
		fmt.Println("DOCDB: " + err.Error())
		return
	}

	err = os.WriteFile("DocDB", data, 0666)
	if err != nil {
		fmt.Println("DOCDB: " + err.Error())
		return
	}
	fmt.Println("Saved DOCDB!")
}

func (d *DocumentDB) ReadDocDB() {

	data, err := os.ReadFile("DocDB")
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
	DocDB.SaveDocDB()
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

	delete(s.Collections, name)
	fmt.Println("Deleted document " + name + " from " + s.Name)
	DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true
}

func (s *Service) GetCollection(name string) *Collection {
	DocDB.M.Lock()
	if _, exist := s.Collections[name]; !exist {
		fmt.Println("Collection with the name " + name + " don't exist")
		DocDB.M.Unlock()
		return nil
	}

	DocDB.M.Unlock()
	return s.Collections[name]
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
	DocDB.SaveDocDB()
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
	DocDB.SaveDocDB()
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

// -------------------------------------------------------------- DOCUMENT FUNCTIONS --------------------------------------------------------------------

func (d *Document) GetContent() interface{} {
	return d.Content
}

func (d *Document) SetContent(name string, content interface{}) bool {

	DocDB.M.Lock()
	d.UpdatedBy = name
	d.Content = content
	DocDB.SaveDocDB()
	DocDB.M.Unlock()
	return true

}
