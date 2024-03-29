package conf

import "testing"

func TestUnmarshalIntoRawMessage(t *testing.T) {
	var d RawMessage
	var data = `
	addr "192.168.10.121"

	log {
		max_size 100
		level "info"
	}

`

	if err := Unmarshal([]byte(data), &d); err != nil {
		t.Fatal(err)
	}

	t.Log(string(d))

}

func TestUnmarshal(t *testing.T) {
	var data = `

age 10 
weight 80  
height 180

`

	var m map[string]int

	if err := Unmarshal([]byte(data), &m); err != nil {
		t.Fatal(err)
	}

	t.Log(m)
}
