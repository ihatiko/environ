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

// ParseCustom parses a custom value based on its type.
// It takes a string value and an r.Value representing the custom type.
// It returns the parsed value as an interface{} and an error if parsing fails.
// The function supports parsing the following types:
// - time.Duration: Parses the string value as a time duration.
// - time.Time: Parses the string value as a time in the format "2006-01-02T15:04:05Z07:00".
// - url.URL: Parses the string value as a URL.
// If the type is not supported, it returns nil and an error of type ErrUnparsableCustom.
func ParseCustom(v string, vValue reflect.Value) (any, error) {
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

func ParsePrimitive(v string, vValue reflect.Value) (any, error) {
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

var splitLexer = "."

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
		paths := strings.Split(name, splitLexer)
		deep(paths, vValue, value, vType)
	}
}

// TODO return error
func deep(paths []string, vValue reflect.Value, value string, vType reflect.Type) reflect.Value {
	for _, path := range paths {
		for i := 0; i < vValue.NumField(); i++ {
			name := vType.Field(i).Name
			if strings.ToUpper(name) != path {
				continue
			}
			field := vValue.Field(i)
			if field.Kind() == reflect.Pointer {
				if !reflect.Indirect(field).IsValid() {
					t := reflect.New(field.Type().Elem())
					vValue.Field(i).Set(t)
					return deep(paths[1:], t.Elem(), value, t.Elem().Type())
				}
				return deep(paths[1:], reflect.Indirect(field), value, vType)

			}
			//TODO err add
			if field.Kind() == reflect.Struct {
				iVType := field.Type()
				parsedValue, err := ParseCustom(value, field)
				if errors.Is(err, ErrUnparsableCustom) {
					deep(paths[1:], field, value, iVType)
					continue
				}
				vValue.Field(i).Set(reflect.ValueOf(parsedValue))
				return vValue
			}
			if field.Kind() == reflect.Map {
				//TODO key exists
				if field.IsNil() {
					field.Set(reflect.MakeMap(field.Type()))
				}
				if len(paths) > 2 {
					newValue := reflect.New(reflect.Indirect(field).Type().Elem())
					deep(paths[2:], newValue.Elem(), value, vType)
					field.SetMapIndex(reflect.ValueOf(strings.ToLower(paths[1])), newValue.Elem())
					return newValue
				}

				parsedValue, err := ParseCustom(value, field)
				if errors.Is(err, ErrUnparsableCustom) {
					parsedPrimitiveValue, err := ParsePrimitive(value, reflect.ValueOf(""))
					if err != nil {
						return field
					}
					if len(paths) > 0 {
						field.SetMapIndex(reflect.ValueOf(strings.ToLower(paths[1])), reflect.ValueOf(parsedPrimitiveValue))
					}
					return field
				}
				field.SetMapIndex(reflect.ValueOf(strings.ToLower(paths[1])), reflect.ValueOf(parsedValue))
				return field
			}
			if field.Kind() == reflect.Slice {
				parsedValue, err := ParseCustom(value, field)
				if errors.Is(err, ErrUnparsableCustom) {
					parsedPrimitiveValue, err := ParsePrimitive(value, reflect.ValueOf(""))
					if err != nil {
						return field
					}
					if len(paths) > 1 {
						index, err := strconv.Atoi(paths[1])
						if err != nil {
							return field
						}
						if field.Len() == 0 {
							field.Set(reflect.MakeSlice(field.Type(), 0, 0))
							if index == 0 {
								slice := reflect.Append(field, reflect.ValueOf(parsedPrimitiveValue))
								field.Set(slice)
							}
						}
						if index+1 > field.Len() {
							diff := index + 1 - field.Len()
							for range diff - 1 {
								slice := reflect.Append(field, reflect.ValueOf(parsedPrimitiveValue))
								field.Set(slice)
							}
						}
					}
					return field
				}
				//TODO кейсы для работы с индексом и работа с кастом типом
				fmt.Println(parsedValue)
				continue
			}
			parsedValue, err := ParseCustom(value, field)
			if errors.Is(err, ErrUnparsableCustom) {
				parsedValue, err = ParsePrimitive(value, field)
				if err != nil {
					return vValue
				}
				resultValue := reflect.ValueOf(parsedValue)
				field.Set(resultValue)
				return resultValue
			}
			if err != nil {
				return vValue
			}
			resultValue := reflect.ValueOf(parsedValue)
			field.Set(resultValue)
			return resultValue
		}
	}
	return vValue
}
