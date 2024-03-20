package environ

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Case2 struct {
	Field2 string
}

type Case1 struct {
	Field1              string
	Field2              float32
	BoolCase            bool
	Field3Duration      time.Duration
	Field4Time          time.Time
	Field3Struct        Case2
	Field3StructPointer *Case2
	Fields              []Case2
	UrlCase             url.URL
}

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		init  func() *Case1
		check func(t *testing.T, value *Case1)
	}{
		{"String check", func() *Case1 {
			os.Setenv("FIELD1", "RESULT")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, "RESULT", value.Field1)
		}},
		{"bool check", func() *Case1 {
			os.Setenv("BoolCase", "true")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, true, value.BoolCase)
		}},
		{"Relplace check", func() *Case1 {
			os.Setenv("FIELD1", "RESULT")
			data := new(Case1)
			data.Field1 = "NotResult"
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, value.Field1, "RESULT")
		}},
		{"Float check", func() *Case1 {
			os.Setenv("FIELD2", "0.1")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, float32(0.1), value.Field2)
		}},
		{"Struct level 1 check", func() *Case1 {
			os.Setenv("FIELD3STRUCT_FIELD2", "RESULT")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, "RESULT", value.Field3Struct.Field2)
		}},
		{"Struct pointer without init level 1 check", func() *Case1 {
			os.Setenv("Field3StructPointer_FIELD2", "RESULT")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, "RESULT", value.Field3StructPointer.Field2)
		}},
		{"Struct pointer with init level 1 check", func() *Case1 {
			os.Setenv("Field3StructPointer_FIELD2", "RESULT")
			data := new(Case1)
			data.Field3StructPointer = new(Case2)
			return data
		}, func(t *testing.T, value *Case1) {
			assert.Equal(t, "RESULT", value.Field3Struct.Field2)
		}},
		{"Check time parse", func() *Case1 {
			os.Setenv("Field4Time", "2006-01-02 15:04:05")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			data, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
			assert.Equal(t, data, value.Field4Time)
		}},
		{"Check duration parse", func() *Case1 {
			os.Setenv("Field3Duration", "10s")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			dur, _ := time.ParseDuration("10s")
			assert.Equal(t, dur, value.Field3Duration)
		}},
		{"Check url parse", func() *Case1 {
			os.Setenv("URLCASE", "http://localhost:10001")
			return new(Case1)
		}, func(t *testing.T, value *Case1) {
			dur, _ := url.Parse("http://localhost:10001")
			assert.Equal(t, *dur, value.UrlCase)
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
