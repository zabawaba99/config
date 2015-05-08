package main

import (
	"log"

	"github.com/zabawaba99/config"
)

type myConfig struct {
	Port   uint   `config:"port"`
	Bucket string `config:"s3_bucket"`
}

func main() {
	var c myConfig
	if err := config.Load(&c); err != nil {
		log.Fatal(err)
	}

	log.Printf("c: %#v\n", c)
}
