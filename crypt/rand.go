package crypt

import (
	"main/config"
	"math/rand"
	"time"
)

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func RandomAESKey() error {
	config.GlobalKey = make([]byte,16)
	_, err := rand.Read(config.GlobalKey[:])
	if err != nil {
		return err
	}
	return nil
}