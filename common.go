package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"time"
)

func ParseJsonFromRequest(req *http.Request, v interface{}) error {
	data := make([]byte, 512)
	count, err := req.Body.Read(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data[:count], v)
}

func UUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		glog.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func EncryptPassword(password string, salt []byte) string {
	if salt == nil {
		m := md5.New()
		m.Write([]byte(time.Now().String()))
		s := hex.EncodeToString(m.Sum(nil))
		salt = []byte(s[2:10])
	}
	mac := hmac.New(sha256.New, salt)
	mac.Write([]byte(password))
	//s := log.Sprintf("%x", (mac.Sum(salt)))
	s := hex.EncodeToString(mac.Sum(nil))

	hasher := sha1.New()
	hasher.Write([]byte(s))

	//result := log.Sprintf("%x", (hasher.Sum(nil)))
	result := hex.EncodeToString(hasher.Sum(nil))

	p := string(salt) + result

	return p
}

func ValidatePassword(hashed string, input_password string) bool {
	salt := hashed[0:8]
	if hashed == EncryptPassword(input_password, []byte(salt)) {
		return true
	} else {
		return false
	}
	return false
}
