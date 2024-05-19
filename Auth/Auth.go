package Auth

import (
	"fmt"
	"math/rand"
)

var KeyDB KeyBase

type KeyBase struct {
	keys     map[string]string
	adminKey string
}

func (k *KeyBase) Init() {
	// Map service name to key
	k.keys = make(map[string]string)
	k.adminKey = "Keks"
}

func (k *KeyBase) CreateAuthenticationKey(name string, keylength int, seed int) (string, bool) {

	if _, exist := k.keys[name]; exist {
		fmt.Println("Key already exists to that name")
		return "", false
	}

	char := []rune("1234567890ABCDEF")
	key := ""

	for i := 1; i < keylength; i++ {
		key += string(char[rand.Intn(seed)%16])
	}

	k.keys[name] = key
	return key, true
}

func (k *KeyBase) CheckAuthenticationKey(s map[string]interface{}) bool {

	if _, ok := s["auth"].(string); !ok {
		fmt.Println("Key type invalid")
		return false
	}

	if _, ok := s["servicename"].(string); !ok {
		fmt.Println("Service name type invalid")
		return false
	}

	if _, exist := k.keys[s["servicename"].(string)]; !exist {
		fmt.Println("Key does not exist")
		return false
	}

	if s["auth"].(string) == k.keys[s["servicename"].(string)] {
		fmt.Println("Key: " + k.keys[s["servicename"].(string)] + " for service: " + s["servicename"].(string) + " is authenticated")
		return true
	}

	fmt.Println("Key did not match")
	return false
}

func (k *KeyBase) CheckAdminKey(s map[string]interface{}) bool {

	if _, ok := s["adminkey"].(string); !ok {
		fmt.Println("Admin key type invalid")
		return false
	}

	if k.adminKey == s["adminkey"].(string) {
		fmt.Println("Admin key match")
		return true
	}

	fmt.Println("Admin key missmatch")
	return false
}
