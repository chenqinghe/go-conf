package conf

import (
	"errors"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	d := &decodeState{}
	d.init(data)

	if err := checkValid(data); err != nil {

	}

	return d.unmarshal(v)
}

func checkValid(data []byte) error {
	// TODO
	return nil
}

type decodeState struct {
	data []byte
	off  int

	items map[string][]byte

	scanner scanner
}

type RawMessage []byte

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("cannot unmarshal into non-pointer or nil")
	}

	d.extractItems()

	switch v.(type) {
	case *RawMessage:
		reflect.Indirect(rv).SetBytes(d.data)
		return nil
	default:
	}

	switch rv.Kind() {
	case reflect.Map:
		d.unmarshalMap(rv)
	case reflect.Struct:
		d.unmarshalStruct(rv)
	}

	return nil
}

func (d *decodeState) unmarshalMap(rv reflect.Value) {
	switch rv.Type().Kind() {
	case reflect.String:
		rv.SetMapIndex(reflect.ValueOf())
	default:
		return errors.New("map key must be string type")
	}
}

func (d *decodeState) unmarshalStruct(rv reflect.Value) {

}

func (d *decodeState) extractItems() {
	d.items = d.scanner.scanKeys()
}

func (d *decodeState) value(v interface{}) error {
	return nil
}

func (d *decodeState) init(data []byte) {
	d.data = data
	d.off = 0
}
