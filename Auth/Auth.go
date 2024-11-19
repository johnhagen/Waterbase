package Auth

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"waterbase/Utils"
)

var KeyDB KeyBase

type KeyBase struct {
	m             sync.Mutex
	keys          map[string]string
	keyLength     int
	adminKey      string
	encryptionKey []byte
}

func (k *KeyBase) Init(adminKey string, keyLength int) {
	k.m.Lock()
	// Map service name to key
	k.keys = make(map[string]string)
	if adminKey == "" {
		log.Fatal("AUTH: Invalid AdminKey Specified")
	}
	k.adminKey = adminKey
	k.keyLength = keyLength
	k.m.Unlock()
}

func (k *KeyBase) CreateAuthenticationKey(name string, keylength int, seed int) (string, bool) {

	k.m.Lock()
	if _, exist := k.keys[name]; exist {
		fmt.Println("Key already exists to that name")
		k.m.Unlock()
		return "", false
	}

	char := []rune("1234567890ABCDEF")
	var key string

	for i := 1; i < keylength; i++ {
		key += string(char[rand.Intn(seed)%len(char)])
	}

	k.keys[name] = key
	k.m.Unlock()
	return key, true
}

func (k *KeyBase) CheckForAuth(s map[string]interface{}) bool {
	k.m.Lock()

	authKeyPresent := Utils.IsString(s["auth"])
	adminKeyPresent := Utils.IsString(s["adminkey"])

	if !authKeyPresent && !adminKeyPresent {
		fmt.Println("No adminkey or authkey specified")
		k.m.Unlock()
		return false
	}

	if adminKeyPresent {
		if s["adminkey"].(string) == k.adminKey {
			k.m.Unlock()
			fmt.Println("Authenticated")
			return true
		} else {
			fmt.Println("Bad Adminkey")
		}
	}

	if authKeyPresent {
		if Utils.IsString(s["servicename"]) {
			if s["auth"].(string) == k.keys[s["servicename"].(string)] {
				k.m.Unlock()
				return true
			}
			k.m.Unlock()
			return false
		}
		k.m.Unlock()
		return false
	}

	k.m.Unlock()
	fmt.Println("Key invalid. Key: " + s["auth"].(string) + " Service: " + s["servicename"].(string))
	return false
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

	if _, exist := k.keys[s["servicename"].(string)]; !exist {
		fmt.Println("Key does not exist")
		k.m.Unlock()
		return false
	}

	if s["auth"].(string) == k.keys[s["servicename"].(string)] {
		//fmt.Println("Key: " + k.Keys[s["servicename"].(string)] + " for service: " + s["servicename"].(string) + " is authenticated")
		k.m.Unlock()
		return true
	}

	k.m.Unlock()
	fmt.Println("Key invalid. Key: " + s["auth"].(string) + " Service: " + s["servicename"].(string))
	return false
}

func (k *KeyBase) CheckAdminKey(s map[string]interface{}) bool {

	k.m.Lock()
	if _, ok := s["adminkey"].(string); !ok {
		//fmt.Println("Admin key type invalid")
		k.m.Unlock()
		return false
	}

	if k.adminKey == s["adminkey"].(string) {
		//fmt.Println("Admin key match")
		k.m.Unlock()
		return true
	}

	fmt.Println("Admin key missmatch")
	k.m.Unlock()
	return false
}

func (k *KeyBase) DeleteKey(name string) bool {
	k.m.Lock()

	if _, exist := k.keys[name]; !exist {
		k.m.Unlock()
		fmt.Printf("Key: %s does not exist\n", name)
		return false
	}

	delete(k.keys, name)
	k.m.Unlock()
	k.SaveDB()
	return true
}

func (k *KeyBase) SaveDB() {
	k.m.Lock()
	data, err := json.Marshal(k.keys)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		k.m.Unlock()
		return
	}

	err = os.WriteFile("KeyDB", data, 0600)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		k.m.Unlock()
		return
	}
	k.m.Unlock()
}

func (k *KeyBase) ReadDB() {
	k.m.Lock()
	data, err := os.ReadFile("KeyDB")
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		k.m.Unlock()
		return
	}

	err = json.Unmarshal(data, &k.keys)
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		k.m.Unlock()
		return
	}
	fmt.Println("Inserted KEYDB File")
	k.m.Unlock()
}

func (k *KeyBase) SaveDB2() {

	data, err := json.Marshal(k.keys)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		return
	}

	encryptedData := Utils.Encrypt(string(k.encryptionKey), string(data))

	err = os.WriteFile("KeyDB", []byte(encryptedData), 0600)
	if err != nil {
		fmt.Println("SAVEKEY: " + err.Error())
		return
	}
}

func (k *KeyBase) ReadDB2() {

	data, err := os.ReadFile("KeyDB")
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		return
	}

	decryptedData := Utils.Decrypt(string(k.encryptionKey), string(data))

	err = json.Unmarshal([]byte(decryptedData), &k.keys)
	if err != nil {
		fmt.Println("READKEY: " + err.Error())
		return
	}
	fmt.Println("Inserted KEYDB File")
}

func (k *KeyBase) encrypt(data string) string {
	encode := Utils.Base64Encode(data)
	return Utils.EncryptAES(k.encryptionKey, string(encode))
}

func (k *KeyBase) decrypt(data string) string {
	decrypt := Utils.DecryptAES(k.encryptionKey, data)
	return string(Utils.Base64Decode(decrypt))
}
