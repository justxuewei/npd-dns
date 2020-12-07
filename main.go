package main

import (
	ndpdns "github.com/xavier-niu/npd-dns/pkg/ndp_dns"
)

func main() {
	dns := ndpdns.NewServer(53)
	dns.LoadConf()
	dns.StartAndServe()
}
