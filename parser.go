package environ

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrUnparsableCustom = errors.New("unparsable custom type")

func ParseCustom(v string, vValue reflect.Value) (interface{}, error) {
	switch vValue.Interface().(type) {
	case time.Duration:
		return time.ParseDuration(v)
	case time.Time:
		return time.Parse(time.DateTime, v)
	case url.URL:
		v, err := url.Parse(v)
		return *v, err
	default:
		return nil, ErrUnparsableCustom
	}
}
func ParsePrimitive(v string, vValue reflect.Value) (interface{}, error) {
	switch vValue.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.String:
		return v, nil
	case reflect.Int:
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	case reflect.Int16:
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	case reflect.Int32:
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	case reflect.Int64:
		return strconv.ParseInt(v, 10, 64)
	case reflect.Int8:
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	case reflect.Uint:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	case reflect.Uint64:
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	case reflect.Uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	case reflect.Float64:
		return strconv.ParseFloat(v, 64)
	case reflect.Float32:
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	default:
		return nil, errors.New("parsed error")
	}
}

func Parse(data any) {
	values := os.Environ()
	vValue := reflect.Indirect(reflect.ValueOf(data))
	vType := vValue.Type()
	for _, i := range values {
		fragments := strings.Split(i, "=")
		if len(fragments) < 2 {
			continue
		}
		name := strings.ToUpper(fragments[0])
		value := fragments[1]
		paths := strings.Split(name, "_")
		deep(paths, vValue, value, vType)
	}
}
func deep(paths []string, vValue reflect.Value, value string, vType reflect.Type) {
	for _, path := range paths {
		for i := 0; i < vValue.NumField(); i++ {
			if strings.ToUpper(vType.Field(i).Name) != path {
				continue
			}
			field := vValue.Field(i)
			if field.Kind() == reflect.Pointer {
				if !reflect.Indirect(field).IsValid() {
					t := reflect.New(field.Type().Elem())
					vValue.Field(i).Set(t)
					deep(paths[1:], t.Elem(), value, t.Elem().Type())
					return
				}
				deep(paths, reflect.Indirect(field), value, vType)
				return
			}
			if field.Kind() == reflect.Struct {
				iVType := field.Type()
				parsedValue, err := ParseCustom(value, field)
				if errors.Is(err, ErrUnparsableCustom) {
					deep(paths[1:], field, value, iVType)
					continue
				}
				vValue.Field(i).Set(reflect.ValueOf(parsedValue))
				return
			}
			if field.Kind() == reflect.Map {
				continue
			}
			if field.Kind() == reflect.Slice {
				if len(paths) > 1 {
					index, err := strconv.Atoi(paths[1])
					if err != nil {
						return
					}
					if field.Len() == 0 {
						fmt.Println(index)
					}
				}

				continue
			}
			parsedValue, err := ParseCustom(value, field)
			if errors.Is(err, ErrUnparsableCustom) {
				parsedValue, err = ParsePrimitive(value, field)
				if err != nil {
					return
				}
				vValue.Field(i).Set(reflect.ValueOf(parsedValue))
			}
			if err != nil {
				return
			}
			vValue.Field(i).Set(reflect.ValueOf(parsedValue))
		}
	}
}
