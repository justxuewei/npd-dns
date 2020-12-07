package main

import (
	ndpdns "github.com/xavier-niu/npd-dns/pkg/ndp_dns"
)

var yamlData = ""

func main() {
	var googleRecords = map[string]string{
		"mail.google.com":  "192.168.0.2",
		"paste.google.com": "192.168.0.3",
	}
	dns := ndpdns.NewServer(53)
	dns.AddZoneData("google.com", googleRecords, nil, ndpdns.DNSForwardLookupZone)
	dns.StartAndServe()
}
