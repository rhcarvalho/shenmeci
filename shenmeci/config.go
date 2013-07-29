package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	Http       *HttpConfig
	CedictPath string
	MongoURL   string
}

type HttpConfig struct {
	Host string
	Port int
}

var config Config

var configFile = flag.String("config", "config.json", "the configuration file in JSON format")

func loadConfig() {
	flag.Parse()
	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&config)
	if err == io.EOF {
		log.Fatal("empty configuration file")
	}
	if err != nil {
		log.Fatal("invalid configuration file: ", err)
	}
}

func validateConfig() {
	errors := []interface{}{}
	if config.Http == nil {
		errors = append(errors, "missing Http configuration")
	} else {
		if len(config.Http.Host) == 0 {
			errors = append(errors, "missing Http.Host configuration")
		}
		if config.Http.Port == 0 {
			errors = append(errors, "missing Http.Port configuration")
		}
	}
	if len(config.CedictPath) == 0 {
		errors = append(errors, "missing CedictPath configuration")
	} else if _, err := os.Stat(config.CedictPath); err != nil {
		errors = append(errors, "invalid CedictPath configuration: ", err)
	}
	if len(config.MongoURL) == 0 {
		errors = append(errors, "missing MongoURL configuration")
	}
	if len(errors) > 0 {
		log.Printf("the configuration file '%s' is invalid\n", *configFile)
		for _, error := range errors {
			log.Println(error)
		}
		fmt.Printf("Example config.json:\n\n%s\n", exampleConfig())
		os.Exit(1)
	}
}

func exampleConfig() []byte {
	b, _ := json.MarshalIndent(
		&Config{&HttpConfig{"127.0.0.1", 8080},
			"dict/cedict_1_0_ts_utf-8_mdbg.txt.gz",
			"localhost:27017"}, "", "  ")
	return b
}
