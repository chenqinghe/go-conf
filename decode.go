package conf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

	switch v.(type) {
	case *RawMessage:
		reflect.Indirect(rv).SetBytes(d.data)
		return nil
	default:
	}

	irv := reflect.Indirect(rv)
	switch irv.Kind() {
	case reflect.Map:
		return d.unmarshalMap(irv)
	case reflect.Struct:
		return d.unmarshalStruct(irv)
	default:
		return errors.New("cannot unmarshal into type" + rv.Type().String())
	}

	return nil
}

func (d *decodeState) unmarshalMap(rv reflect.Value) error {
	// check for map key type
	switch rv.Type().Key().Kind() {
	case reflect.String:
	default:
		return errors.New("map key must be string type")
	}
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	d.extractItems()
	fmt.Println(d.items)

	switch kind := rv.Type().Elem().Kind(); kind {
	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
		reflect.Int, reflect.Uint:
		for key, val := range d.items {
			fmt.Println(key, string(val))
			d.scanner.init(val)
			d.scanner.skipWhitespace()
			tk, lit := d.scanner.scanNumber(false)
			fmt.Println("tk:", tk, "lit:", lit)
			var v interface{}
			switch tk {
			case INT:
				v, _ = strconv.ParseInt(lit, 10, 64)
			default:
				fmt.Println("type error")
				return errors.New("cannot unmarshal " + tokens[tk] + " into type " + kind.String())
			}
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(v).Convert(rv.Type().Elem()))
		}
	case reflect.Float32, reflect.Float64:
		for key, val := range d.items {
			d.scanner.init(val)
			tk, lit := d.scanner.scanNumber(false)
			var v interface{}
			switch tk {
			case FLOAT:
				v, _ = strconv.ParseFloat(lit, 64)
			default:
				return errors.New("cannot unmarshal " + tokens[tk] + " into type " + kind.String())
			}
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(v).Convert(rv.Elem().Type()))
		}
	case reflect.String:
	case reflect.Array, reflect.Slice:
	case reflect.Struct:
	}

	return nil
}

func (d *decodeState) unmarshalStruct(rv reflect.Value) error {
	return nil
}

func (d *decodeState) unmarshalArray(rv reflect.Value) error {
	return nil
}

func (d *decodeState) unmarshalLiteral(rv reflect.Value) error {
	return nil
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
	d.scanner.init(data)
}

func (d *decodeState) int(rv reflect.Value) {

}
