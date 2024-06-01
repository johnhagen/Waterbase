package DocumentDB

import (
	"encoding/json"
	"fmt"
	"os"
)

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
