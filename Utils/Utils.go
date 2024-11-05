package Utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ReadHeader(r *http.Request) map[string]interface{} {

	//fmt.Println(r.Header)

	data := make(map[string]interface{})

	data["adminkey"] = r.Header.Get("Adminkey")
	data["auth"] = r.Header.Get("Auth")
	data["servicename"] = r.Header.Get("Servicename")
	data["collectionname"] = r.Header.Get("Collectionname")
	data["documentname"] = r.Header.Get("Documentname")

	return data
}

func ReadFromJSON(r *http.Request) (map[string]interface{}, error) {

	data := make(map[string]interface{})

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func IsString(s interface{}) bool {
	if _, ok := s.(string); ok {
		return true
	}
	return false
}

func Encrypt(keyString string, stringToEncrypt string) (encryptedString string) {
	// convert key to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func Decrypt(keyString string, stringToDecrypt string) string {
	key, _ := hex.DecodeString(keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func EncryptAES(key []byte, text string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	out := make([]byte, len(text))

	c.Encrypt(out, []byte(text))

	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return string(pt[:])
}

func Base64Encode(s string) []byte {
	return []byte(base64.StdEncoding.EncodeToString([]byte(s)))
}

func Base64Decode(s string) []byte {
	decode, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return decode
}
