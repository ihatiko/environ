package environ

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultStringResult = "RESULT"
)

type Level1 struct {
	FieldString string
}
type Level2 struct {
	FieldStruct        Level1
	FieldStructPointer *Level1
}
type Case1 struct {
	FieldString             string
	FieldFloat32            float32
	FieldBool               bool
	FieldDuration           time.Duration
	FieldTime               time.Time
	FieldStruct             Level1
	FieldStructPointer      *Level1
	Fields                  []Level1
	FieldUrl                url.URL
	FieldMap                map[string]string
	FieldMapStruct          map[string]Level1
	FieldInnerStructPointer *Level2
	FieldInnerStruct        Level2
}

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		init  func() *Case1
		check func(t *testing.T, value *Case1)
	}{
		{"Struct inner pointer without init level 2 check", func() *Case1 {
			os.Setenv("FieldInnerStructPointer_FieldStructPointer_FieldString", defaultStringResult)
			data := new(Case1)
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldInnerStructPointer.FieldStructPointer.FieldString)
		}},
		{"Struct inner with init level 2 check", func() *Case1 {
			os.Setenv("FieldInnerStruct_FieldStruct_FieldString", defaultStringResult)
			data := new(Case1)
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldInnerStruct.FieldStruct.FieldString)
		}},
		{"String check", func() *Case1 {
			os.Setenv("FieldString", defaultStringResult)
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldString)
		}},
		{"bool check", func() *Case1 {
			os.Setenv("FieldBool", "true")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, true, value.FieldBool)
		}},
		{"Replace check", func() *Case1 {
			os.Setenv("FieldString", defaultStringResult)
			data := new(Case1)
			data.FieldString = "NotResult"
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, value.FieldString, defaultStringResult)
		}},
		{"Float check", func() *Case1 {
			os.Setenv("FieldFloat32", "0.1")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, float32(0.1), value.FieldFloat32)
		}},
		{"Struct level 1 check", func() *Case1 {
			os.Setenv("FieldStruct_FieldString", defaultStringResult)
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldStruct.FieldString)
		}},
		{"Struct pointer without init level 1 check", func() *Case1 {
			os.Setenv("FieldStructPointer_FieldString", defaultStringResult)
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldStructPointer.FieldString)
		}},
		{"Struct pointer with init level 1 check", func() *Case1 {
			os.Setenv("FieldStructPointer_FieldString", defaultStringResult)
			data := new(Case1)
			data.FieldStructPointer = new(Level1)
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, defaultStringResult, value.FieldStructPointer.FieldString)
		}},
		{"Check time parse", func() *Case1 {
			os.Setenv("FieldTime", "2006-01-02 15:04:05")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			data, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
			assert.Equal(t, data, value.FieldTime)
		}},
		{"Check duration parse", func() *Case1 {
			os.Setenv("FieldDuration", "10s")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			dur, _ := time.ParseDuration("10s")
			assert.Equal(t, dur, value.FieldDuration)
		}},
		{"Check url parse", func() *Case1 {
			os.Setenv("FieldUrl", "http://localhost:10001")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			dur, _ := url.Parse("http://localhost:10001")
			assert.Equal(t, *dur, value.FieldUrl)
		}},
		{"Check map parse with nil", func() *Case1 {
			os.Setenv("FieldMap_Key", defaultStringResult)
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, value.FieldMap["key"], defaultStringResult)
		}},
		{"Check map parse", func() *Case1 {
			os.Setenv("FieldMap_Key", defaultStringResult)
			case1 := new(Case1)
			case1.FieldMap = make(map[string]string)
			return case1
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, value.FieldMap["key"], defaultStringResult)
		}},
		{"Check map parse struct", func() *Case1 {
			os.Setenv("FieldMapStruct_Key_FieldString", defaultStringResult)
			case1 := new(Case1)
			case1.FieldMapStruct = make(map[string]Level1)
			return case1
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, value.FieldMapStruct["key"].FieldString, defaultStringResult)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.init()
			Parse(data)
			tt.check(t, data)
		})
	}
}
