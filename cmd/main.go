package main

import (
	"findApi/bootstrap"
	"findApi/internal/encryptutil"
	"log"
)

func main(){
	env := bootstrap.LoadEnv()
	log.Println(env)
	enc := encryptutil.Encrypt("hello",encryptutil.EncryptionKey)
	log.Println(enc)
	dec := encryptutil.Decrypt(enc,encryptutil.EncryptionKey)
	log.Println(dec)
}