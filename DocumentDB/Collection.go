package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
	CacheMem "waterbase/Cache"
	consts "waterbase/Data"
)

// -------------------------------------------------------------- COLLECTION FUNCTIONS --------------------------------------------------------------------

func (c *Collection) CreateNewDocument(name string, owner string, content interface{}) bool {

	DocDB.M.Lock()
	_, err := os.Stat(consts.DEFAULT_SAVE_LOCATION + c.ServiceName + "/" + c.Name + "/" + name)
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

	document.SaveDocument(consts.DEFAULT_SAVE_LOCATION + c.ServiceName + "/" + c.Name)
	fmt.Println("Service: " + c.ServiceName + " - Created document: " + name)
	DocDB.M.Unlock()
	return true
}

func (c *Collection) GetDocument(name string) *Document {
	DocDB.M.Lock()

	cachedData := CacheMem.Cache.Get("doc-" + name)
	if cachedData == nil {
		file, err := os.ReadFile(consts.DEFAULT_SAVE_LOCATION + c.ServiceName + "/" + c.Name + "/" + name)
		if err != nil {
			fmt.Println(err.Error())
			DocDB.M.Unlock()
			return nil
		}
		cachedData = &file
		CacheMem.Cache.Insert("doc-"+name, file)
	}

	document := Document{}

	err := json.Unmarshal(*cachedData, &document)
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
	err := os.Remove(consts.DEFAULT_SAVE_LOCATION + c.ServiceName + "/" + c.Name + "/" + name)
	if err != nil {
		fmt.Println(err.Error())
		DocDB.M.Unlock()
		return false
	}
	CacheMem.Cache.Delete("doc-" + name)
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
