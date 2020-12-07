package ndp_dns

import (
	"testing"
)

func TestConfRead(t *testing.T) {
	//_ = os.Setenv("NPD_DNS_CONF", "/Users/xavier/Development/Go/npd-dns/conf.yaml")
	conf := Conf{}
	records, err := conf.Read()
	if err != nil {
		t.Error(err)
	}
	t.Log(records)
}
