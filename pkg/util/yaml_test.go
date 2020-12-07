package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	records, err := read("./test.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	println(records)
}

func TestWrite(t *testing.T) {
	var googleRecords = map[string]string{
		"mail.google.com":  "192.168.0.2",
		"paste.google.com": "192.168.0.3",
	}
	var firstMap = map[string]map[string]string{
		"google.com": googleRecords,
	}
	d, err := yaml.Marshal(&firstMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fileName := "test.yaml"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	s := string(d)
	dstFile.WriteString(s + "\n")
}
