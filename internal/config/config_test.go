package config

import (
	"reflect"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

type testConfig struct {
	MyStr    string `env:"MY_STR"`
	MyInt    int    `env:"MY_INT"`
	MyUint   uint   `env:"MY_UINT"`
	MyBool   bool   `env:"MY_BOOL"`
	MyStruct struct {
		Nested1        string `env:"NESTED_1"`
		MyNestedStruct struct {
			Nested2 string `env:"NESTED_2"`
		}
	}
}

func TestPopulate(t *testing.T) {
	t.Setenv("MY_STR", "str1")
	t.Setenv("MY_INT", "123")
	t.Setenv("NESTED_1", "nested1")
	t.Setenv("NESTED_2", "nested2")
	cfg := &testConfig{}

	assert.NoError(t, populate(reflect.ValueOf(cfg)))

	assert.Equal(t, "str1", cfg.MyStr)
	assert.Equal(t, 123, cfg.MyInt)
	assert.Equal(t, uint(0), cfg.MyUint) // env not set
	assert.Equal(t, false, cfg.MyBool)   // env not set
	assert.Equal(t, "nested1", cfg.MyStruct.Nested1)
	assert.Equal(t, "nested2", cfg.MyStruct.MyNestedStruct.Nested2)

	t.Setenv("MY_UINT", "124")
	t.Setenv("MY_BOOL", "true")

	cfg = &testConfig{}

	assert.NoError(t, populate(reflect.ValueOf(cfg)))

	assert.Equal(t, "str1", cfg.MyStr)
	assert.Equal(t, 123, cfg.MyInt)
	assert.Equal(t, uint(124), cfg.MyUint)
	assert.Equal(t, true, cfg.MyBool)
	assert.Equal(t, "nested1", cfg.MyStruct.Nested1)
	assert.Equal(t, "nested2", cfg.MyStruct.MyNestedStruct.Nested2)
}
