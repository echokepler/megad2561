package qsparser

import (
	"github.com/echokepler/megad2561/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Enum int

const (
	One Enum = iota
	Two
	Three
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		t        interface{}
		expected core.ServiceValues
	}{
		{
			name: "Should be correct marshal integer values",
			t: struct {
				Uint   uint   `qs:"uint"`
				Uint8  uint8  `qs:"uint8"`
				Uint16 uint16 `qs:"uint16"`
				Uint32 uint32 `qs:"uint32"`
				Uint64 uint64 `qs:"uint64"`
				Int    int    `qs:"int"`
				Int8   int8   `qs:"int8"`
				Int16  int16  `qs:"int16"`
				Int32  int32  `qs:"int32"`
				Int64  int64  `qs:"int64"`
			}{
				Uint:   1,
				Uint8:  100,
				Uint16: 1000,
				Uint32: 10000,
				Uint64: 100000,
				Int:    1,
				Int8:   100,
				Int16:  1000,
				Int32:  10000,
				Int64:  100000,
			},
			expected: core.ServiceValues{
				"uint":   []string{"1"},
				"uint8":  []string{"100"},
				"uint16": []string{"1000"},
				"uint32": []string{"10000"},
				"uint64": []string{"100000"},
				"int":    []string{"1"},
				"int8":   []string{"100"},
				"int16":  []string{"1000"},
				"int32":  []string{"10000"},
				"int64":  []string{"100000"},
			},
		},
		{
			name: "Should be correct marshal boolean values",
			t: struct {
				Bool bool `qs:"b"`
			}{
				Bool: false,
			},
			expected: core.ServiceValues{
				"b": []string{"false"},
			},
		},
		{
			name: "Should be correct marshal strings values",
			t: struct {
				String string `qs:"str"`
			}{
				String: "some string",
			},
			expected: core.ServiceValues{
				"str": []string{"some string"},
			},
		},
		{
			name: "Should be correct marshal enum values",
			t: struct {
				SomeEnum Enum `qs:"enum"`
			}{
				SomeEnum: Two,
			},
			expected: core.ServiceValues{
				"enum": []string{"1"},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualValues(t, tc.expected, Marshal(tc.t))
		})
	}

	t.Run("Should skip empty tag", func(t *testing.T) {
		tCase := struct {
			Name      string `qs:"name"`
			Nickname  string `qs:"nick"`
			HiddenAge uint8
			Age       uint8 `qs:"age"`
		}{
			Name:      "Anton",
			Nickname:  "devdammit",
			HiddenAge: 24,
			Age:       24,
		}

		values := Marshal(tCase)
		expected := core.ServiceValues{
			"name": []string{tCase.Name},
			"nick": []string{tCase.Nickname},
			"age":  []string{"24"},
		}

		assert.EqualValues(t, expected, values)
	})
}

func TestUnMarshal(t *testing.T) {
	t.Parallel()

	type Integers struct {
		Uint   uint   `qs:"uint"`
		Uint8  uint8  `qs:"uint8"`
		Uint16 uint16 `qs:"uint16"`
		Uint32 uint32 `qs:"uint32"`
		Uint64 uint64 `qs:"uint64"`
		Int    int    `qs:"int"`
		Int8   int8   `qs:"int8"`
		Int16  int16  `qs:"int16"`
		Int32  int32  `qs:"int32"`
		Int64  int64  `qs:"int64"`
	}

	t.Run("Should be correct unmarshal integers", func(t *testing.T) {
		var integers Integers
		values := core.ServiceValues{
			"uint":   []string{"1"},
			"uint8":  []string{"100"},
			"uint16": []string{"1000"},
			"uint32": []string{"10000"},
			"uint64": []string{"100000"},
			"int":    []string{"1"},
			"int8":   []string{"100"},
			"int16":  []string{"1000"},
			"int32":  []string{"10000"},
			"int64":  []string{"100000"},
		}

		err := UnMarshal(values, &integers)
		if err != nil {
			t.Error(err)
		}

		assert.NotEmpty(t, integers)
		assert.EqualValues(t, Integers{
			Uint:   1,
			Uint8:  100,
			Uint16: 1000,
			Uint32: 10000,
			Uint64: 100000,
			Int:    1,
			Int8:   100,
			Int16:  1000,
			Int32:  10000,
			Int64:  100000,
		}, integers)
	})
}
