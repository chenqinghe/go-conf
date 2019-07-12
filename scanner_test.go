package conf

import (
	"encoding/json"
	"testing"
)

func TestScanKeys(t *testing.T) {
	var data = `
addr "192.168.107.246"

log {
	level "info"
	path "/data/logs/conf.log"
}

es {
	index "test"
	cluster {
		addr ["192.168.10.151"]
		balance round
	}
}

`

	var s = &scanner{}
	s.init([]byte(data))

	for k, v := range s.scanKeys() {
		t.Log(k, string(v))
	}

}

func TestUnmarshalToMap(t *testing.T) {
	var data = `
{
	"addr":"192.168.107.246",
	"log":{
			"level": "info",
			"path":"/data/logs/conf.log"
	},
	"es":{
		"index": "test",
		"cluster":{
			"addr":["192.168.10.151"],
			"balance":"round" 
		}
	}
}
`


	m := map[string]string{}

	if err := json.Unmarshal([]byte(data), &m); err != nil {
		t.Fatal(err)
	}

	t.Log(m)
}
