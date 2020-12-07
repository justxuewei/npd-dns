package ndp_dns

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Env string

const (
	ConfPath Env = "NPD_DNS_CONF"
)

type RecordMap map[string]string

type Records struct {
	DomainName string
	Map        RecordMap
}

type Conf struct {
	path string
}

func (c *Conf) Read() ([]Records, error) {
	if c.path == "" {
		c.path = os.Getenv(string(ConfPath))
		if c.path == "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return nil, err
			}
			c.path = dir + "/conf.yaml"
		}
	}
	conf, err := readConf(c.path)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func readConf(path string) ([]Records, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	yamlData := string(content)
	rootMap := make(map[string]map[string]map[string]string)
	err = yaml.Unmarshal([]byte(yamlData), &rootMap)
	if err != nil {
		return nil, err
	}
	firstMap := rootMap["records"]
	records := make([]Records, 0)
	for k, v := range firstMap {
		records = append(records, Records{
			DomainName: k,
			Map:        v,
		})
	}
	return records, nil
}
