package main

import (
	ndpdns "github.com/xavier-niu/npd-dns/pkg/ndp_dns"
	"gopkg.in/yaml.v2"
	"log"
)

var yamlData = ""

func main() {
	var googleRecords = map[string]string{
		"mail.google.com":  "192.168.0.2",
		"paste.google.com": "192.168.0.3",
	}
	firstMap := make(map[string]map[string]string)
	err := yaml.Unmarshal([]byte(yamlData), &firstMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	dns := ndpdns.NewServer(10001)
	dns.AddZoneData("google.com", googleRecords, nil, ndpdns.DNSForwardLookupZone)
	dns.StartAndServe()
}
