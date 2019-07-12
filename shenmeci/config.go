package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Http       *HttpConfig
	StaticPath string
	CedictPath string
}

type HttpConfig struct {
	Host string
	Port int
}

var config Config

var configFile = flag.String("config", "config.json", "the configuration file in JSON format")

func loadConfig() {
	flag.Parse()
	configFileAbsPath, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	configFile = &configFileAbsPath
	file, err := os.Open(*configFile)
	if err != nil {
		logFatalAndExampleConfig(err)
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&config)
	if err == io.EOF {
		logFatalAndExampleConfig(fmt.Sprint("empty configuration file: ", *configFile))
	}
	if err != nil {
		logFatalAndExampleConfig(fmt.Sprintf("the configuration file '%s' is invalid: %v", *configFile, err))
	}
	absRelToConfigFile := func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(filepath.Dir(*configFile), path)
	}
	config.StaticPath = absRelToConfigFile(config.StaticPath)
	config.CedictPath = absRelToConfigFile(config.CedictPath)
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
	if len(config.StaticPath) == 0 {
		errors = append(errors, "missing StaticPath configuration")
	} else {
		fi, err := os.Stat(config.StaticPath)
		if err != nil {
			errors = append(errors, fmt.Sprint("invalid StaticPath configuration: ", err))
		} else if !fi.IsDir() {
			errors = append(errors, fmt.Sprint("StaticPath should be a directory: ", config.StaticPath))
		}
	}
	if len(config.CedictPath) == 0 {
		errors = append(errors, "missing CedictPath configuration")
	} else if _, err := os.Stat(config.CedictPath); err != nil {
		errors = append(errors, fmt.Sprint("invalid CedictPath configuration: ", err))
	}
	if len(errors) > 0 {
		errors = append([]interface{}{fmt.Sprintf("the configuration file '%s' is invalid", *configFile)},
			errors...)
		logFatalAndExampleConfig(errors...)
	}
}

func exampleConfig() []byte {
	b, _ := json.MarshalIndent(
		&Config{&HttpConfig{"127.0.0.1", 8080},
			"static/",
			"dict/cedict_1_0_ts_utf-8_mdbg.txt.gz",
		}, "", "  ")
	return b
}

func logFatalAndExampleConfig(error ...interface{}) {
	for _, err := range error {
		log.Println(err)
	}
	fmt.Printf("Example config.json:\n\n%s\n", exampleConfig())
	os.Exit(1)
}
