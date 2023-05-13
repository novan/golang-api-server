package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	jsoniter "github.com/json-iterator/go"
)

var (
	Json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func AtoI(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func AtoI64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func AtoF32(s string) float32 {
	return float32(AtoF64(s))
}

func AtoF64(s string) float64 {
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return f
	}
	return 0
}

// GenerateToken returns a unique token based on the provided email string
func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateVerificationCode() int {
	_min := 1000
	_max := 9999
	return rand.Intn((_max - _min) + _min)
}

type MessageData struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func DecodeMessageData(b []byte, i interface{}) error {
	md := &MessageData{Data: i}
	if err := json.Unmarshal(b, &md); err != nil {
		return err
	}
	return nil
}

func ToString(res interface{}) string {
	resStr, _ := json.Marshal(res)
	return string(resStr)
}

func IsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

type RawTime []byte

func (t RawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

func ShortStack(stack []byte) string {
	var stackStr string
	needle := "goroutine"
	stackStr = string(stack[:])
	cut := strings.Index(stackStr[1:], needle)
	if cut > 0 {
		return stackStr[:cut]
	}
	return stackStr
}

func IfNull(new interface{}, def interface{}) interface{} {
	if new == nil || reflect.ValueOf(new).IsNil() {
		return def
	}
	return new
}

func MakeUniqueFilename(fileName string) (string, error) {
	nameOnly := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	ext := filepath.Ext(fileName)

	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s-%x%s", nameOnly, b, ext)
	return s, nil
}