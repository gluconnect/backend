package models

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"time"
)

const (
	BeforeMeal byte = iota
	AfterMeal  byte = iota
)

type GlucoseReading struct {
	Timestamp     time.Time `json:"timestamp"`
	Reading       float64   `json:"reading"`
	ReadingType   byte      `json:"readingtype"`
	MeasureMethod string    `json:"measuremethod"`
}

type User struct {
	Email       string            `json:"email"`
	Password    [sha512.Size]byte `json:"password"`
	LinkedUsers []string          `json:"linked_users"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LinkUserRequest struct {
	Email string `json:"link"`
}

func DecodeStruct[K any](buf []byte) K {
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	var ret K
	dec.Decode(&ret)
	return ret
}

func EncodeStruct[K any](user K) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(user)
	return buf.Bytes()
}
