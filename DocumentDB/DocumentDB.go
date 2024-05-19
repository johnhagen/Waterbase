package DocumentDB

import "fmt"

var DocDB DocumentDB

type DocumentDB struct {
	Services map[string]*Service
}

type Service struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Collections map[string]*Collection
}

type Collection struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	LastUpdated string
	Documents   map[string]*Document
}

type Document struct {
	UpdatedBy    string
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	CreationDate string
	LastUpdated  string
	Content      interface{} `json:"content"`
}

// -------------------------------------------------------------- DB FUNCTIONS --------------------------------------------------------------------

func (d *DocumentDB) InitDB() {
	d.Services = make(map[string]*Service)
}

func (d *DocumentDB) CreateNewService(r Service) bool {

	if _, exist := d.Services[r.Name]; exist {
		fmt.Println("Service with the name " + r.Name + " already exists")
		return false
	}

	service := Service{}

	service.Collections = make(map[string]*Collection)
	service.Name = r.Name
	service.Owner = r.Owner

	d.Services[r.Name] = &service
	fmt.Println("Created service " + r.Name)
	return true
}

func (d *DocumentDB) GetService(name string) *Service {
	if _, exist := d.Services[name]; !exist {
		fmt.Println("Could not find service " + name)
		return nil
	}

	return d.Services[name]
}

// -------------------------------------------------------------- SERVICE FUNCTIONS --------------------------------------------------------------------

func (s *Service) CreateNewCollection(name string, owner string) bool {

	if _, exist := s.Collections[name]; exist {
		fmt.Println("Collection with the name " + name + " already exists")
		return false
	}

	var collection Collection

	collection.Documents = make(map[string]*Document)
	collection.Name = name
	collection.Owner = owner

	s.Collections[name] = &collection
	fmt.Println("Created collection " + name)
	return true
}

func (s *Service) DeleteCollection(name string) bool {

	if _, exist := s.Collections[name]; !exist {
		fmt.Println("Collection with the name " + name + " don't exist")
		return false
	}

	delete(s.Collections, name)
	fmt.Println("Deleted document " + name + " from " + s.Name)
	return true
}

func (s *Service) GetCollection(name string) *Collection {

	if _, exist := s.Collections[name]; !exist {
		fmt.Println("Collection with the name " + name + " don't exist")
		return nil
	}

	return s.Collections[name]
}

// -------------------------------------------------------------- COLLECTION FUNCTIONS --------------------------------------------------------------------

func (c *Collection) CreateNewDocument(name string, owner string, content interface{}) bool {

	if _, exist := c.Documents[name]; exist {
		fmt.Println("document with the name " + name + " already exists")
		return false
	}

	var document Document

	document.Name = name
	document.Owner = owner
	document.UpdatedBy = owner
	document.Content = content

	c.Documents[name] = &document
	fmt.Println("Created document" + name)
	return true
}

func (c *Collection) DeleteDocument(name string) bool {

	if _, exist := c.Documents[name]; !exist {
		fmt.Println("document with the name " + name + " don't exist")
		return false
	}

	delete(c.Documents, name)
	fmt.Println("Deleted document " + name + " from " + c.Name)
	return true
}

func (c *Collection) GetDocument(name string) *Document {

	if _, exist := c.Documents[name]; !exist {
		fmt.Println("document with the name " + name + " don't exist")
		return nil
	}

	return c.Documents[name]
}

// -------------------------------------------------------------- DOCUMENT FUNCTIONS --------------------------------------------------------------------

func (d *Document) GetContent() interface{} {

	return d.Content

}

func (d *Document) SetContent(name string, content interface{}) bool {

	d.UpdatedBy = name
	d.Content = content
	return true

}
