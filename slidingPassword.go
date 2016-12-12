package main // github.com/gitinsky/vnc-go-web

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/satori/go.uuid"
	"regexp"
	"sync"
	"time"
)

type SlidingPassword struct {
	curPass []byte
	oldPass []byte
	iv      []byte
	lock    sync.RWMutex
}

func NewSlidingPassword() *SlidingPassword {
	return &SlidingPassword{curPass: uuidgen(), oldPass: uuidgen(), iv: uuidgen()}
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
	p.curPass = uuidgen()
}

func (p *SlidingPassword) GetPasswords() ([]byte, []byte, []byte) {
	defer p.lock.RUnlock()
	p.lock.RLock()
	return p.curPass, p.oldPass, p.iv
}

func (p *SlidingPassword) Decrypt(msg []byte, check *regexp.Regexp) []string {
	curPass, oldPass, iv := p.GetPasswords()

	res := check.FindAllStringSubmatch(DecryptAESCFB(msg, curPass, iv), -1)
	if res != nil {
		return res[0]
	}

	res = check.FindAllStringSubmatch(DecryptAESCFB([]byte(msg), oldPass, iv), -1)
	if res != nil {
		return res[0]
	}

	return nil
}

func (p *SlidingPassword) Encrypt(msg string) []byte {
	curPass, _, iv := p.GetPasswords()

	return EncryptAESCFB([]byte(msg), curPass, iv)
}

func uuidgen() []byte {
	data, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		panic(err)
	}

	return data
}

func EncryptAESCFB(src, key, iv []byte) []byte {
	dst := make([]byte, len(src))

	c := cipher.NewCFBEncrypter(aesCipher(key), iv)
	c.XORKeyStream(dst, src)

	//cipher.NewCFBEncrypter(aesCipher(key), iv).XORKeyStream(dst, src)

	return dst
}

func DecryptAESCFB(src, key, iv []byte) string {
	dst := make([]byte, len(src))

	cipher.NewCFBDecrypter(aesCipher(key), iv).XORKeyStream(dst, src)

	return string(dst)
}

func aesCipher(key []byte) cipher.Block {
	cip, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cip
}
