package Auth

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

var KeyDB KeyBase

type KeyBase struct {
	m         sync.Mutex
	Keys      map[string]string `json:"keys"`
	KeyLength int               `json:"keyLength"`
	AdminKey  string            `json:"adminKey"`
}

func (k *KeyBase) Init(adminKey string, keyLength int) {
	// Map service name to key
	k.Keys = make(map[string]string)
	k.AdminKey = adminKey
	k.KeyLength = keyLength
}

func (k *KeyBase) CreateAuthenticationKey(name string, keylength int, seed int) (string, bool) {

	k.m.Lock()
	if _, exist := k.Keys[name]; exist {
		fmt.Println("Key already exists to that name")
		k.m.Unlock()
		return "", false
	}

	char := []rune("1234567890ABCDEF")
	key := ""

	for i := 1; i < keylength; i++ {
		key += string(char[rand.Intn(seed)%16])
	}

	k.Keys[name] = key
	k.m.Unlock()
	return key, true
}

func (k *KeyBase) CheckAuthenticationKey(s map[string]interface{}) bool {

	k.m.Lock()
	if _, ok := s["auth"].(string); !ok {
		fmt.Println("Key type invalid")
		k.m.Unlock()
		return false
	}

	if _, ok := s["servicename"].(string); !ok {
		fmt.Println("Service name type invalid")
		k.m.Unlock()
		return false
	}

	if _, exist := k.Keys[s["servicename"].(string)]; !exist {
		fmt.Println("Key does not exist")
		k.m.Unlock()
		return false
	}

	if s["auth"].(string) == k.Keys[s["servicename"].(string)] {
		fmt.Println("Key: " + k.Keys[s["servicename"].(string)] + " for service: " + s["servicename"].(string) + " is authenticated")
		k.m.Unlock()
		return true
	}

	k.m.Unlock()
	fmt.Println("Key did not match")
	return false
}

func (k *KeyBase) CheckAdminKey(s map[string]interface{}) bool {

	k.m.Lock()
	if _, ok := s["adminkey"].(string); !ok {
		fmt.Println("Admin key type invalid")
		k.m.Unlock()
		return false
	}

	if k.AdminKey == s["adminkey"].(string) {
		fmt.Println("Admin key match")
		k.m.Unlock()
		return true
	}

	fmt.Println("Admin key missmatch")
	k.m.Unlock()
	return false
}

func (k *KeyBase) SaveDB() {

	data, err := json.Marshal(k.Keys)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		return
	}

	err = os.WriteFile("KeyDB", data, 0666)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		return
	}
}

func (k *KeyBase) ReadDB() {
	data, err := os.ReadFile("KeyDB")
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		return
	}

	err = json.Unmarshal(data, &k.Keys)
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		return
	}
	fmt.Println("Inserted KEYDB File")
}
