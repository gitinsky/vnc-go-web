package main // github.com/gitinsky/vnc-go-web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"hash/crc32"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

type EncrData struct {
	IV   []byte
	Data []byte
}

type RawData struct {
	CSum uint32
	Data []byte
}

const (
	passwdLen = 16
)

func randStr(strLen int) []byte {
	buf := make([]byte, strLen)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

type SlidingPassword struct {
	curPass []byte
	oldPass []byte
	lock    sync.RWMutex
}

func NewSlidingPassword() *SlidingPassword {
	return &SlidingPassword{curPass: randStr(passwdLen), oldPass: randStr(passwdLen)}
}

func (p *SlidingPassword) UpdateLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		p.Update()
	}
}

func (p *SlidingPassword) Update() {
	defer p.lock.Unlock()
	p.lock.Lock()
	p.oldPass = p.curPass
	p.curPass = randStr(passwdLen)
}

func (p *SlidingPassword) GetPasswords() ([]byte, []byte) {
	defer p.lock.RUnlock()
	p.lock.RLock()
	return p.curPass, p.oldPass
}

func (p *SlidingPassword) Decrypt(msg []byte) *AuthToken {
	curPass, oldPass := p.GetPasswords()

	res := DecryptAESCFB(msg, curPass)
	if res != nil {
		return res
	}

	res = DecryptAESCFB(msg, oldPass)
	if res != nil {
		return res
	}

	return nil
}

func (p *SlidingPassword) Encrypt(authToken *AuthToken) []byte {
	curPass, _ := p.GetPasswords()

	return EncryptAESCFB(authToken, curPass)
}

func uuidgen() []byte {
	data, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		panic(err)
	}

	return data
}

func EncryptAESCFB(src *AuthToken, key []byte) []byte {
	var err error
	rawData := RawData{}

	rawData.Data, err = json.Marshal(src)
	if err != nil {
		panic(err)
	}
	rawData.CSum = crc32.ChecksumIEEE(rawData.Data)

	rawBytes, err := json.Marshal(rawData)
	encrData := EncrData{
		IV:   randStr(len(key)),
		Data: make([]byte, len(rawBytes)),
	}

	cipher.NewCFBEncrypter(aesCipher(key), encrData.IV).XORKeyStream(encrData.Data, rawBytes)

	encrBytes, err := json.Marshal(encrData)
	if err != nil {
		panic(err)
	}

	return encrBytes
}

func DecryptAESCFB(src, key []byte) *AuthToken {
	encrData := EncrData{}
	err := json.Unmarshal(src, &encrData)
	if err != nil {
		return nil
	}

	if len(encrData.IV) != len(key) {
		return nil
	}

	rawBytes := make([]byte, len(encrData.Data))
	cipher.NewCFBDecrypter(aesCipher(key), encrData.IV).XORKeyStream(rawBytes, encrData.Data)

	rawData := RawData{}
	err = json.Unmarshal(rawBytes, &rawData)
	if err != nil {
		return nil
	}

	cSum := crc32.ChecksumIEEE(rawData.Data)
	if cSum != rawData.CSum {
		return nil
	}

	res := AuthToken{}
	err = json.Unmarshal(rawData.Data, &res)
	if err != nil {
		return nil
	}

	return &res
}

func aesCipher(key []byte) cipher.Block {
	cip, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cip
}
