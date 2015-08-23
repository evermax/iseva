package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseRandom(t *testing.T) {
	tests := []struct {
		params   funcParams
		size     int
		key      string
		foundKey bool
	}{
		{
			params: funcParams{
				Randoms: map[string]random{
					"string": random{
						Type: "string",
						Max:  0,
						Min:  0,
						Size: 20,
					},
				},
			},
			size:     1,
			key:      "string",
			foundKey: true,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"int": random{
						Type: "int",
						Max:  230,
						Min:  -22,
					},
				},
			},
			size:     1,
			key:      "int",
			foundKey: true,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
				},
			},
			size:     1,
			key:      "float",
			foundKey: true,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
				},
			},
			size:     1,
			key:      "flat",
			foundKey: false,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
					"int": random{
						Type: "int",
						Max:  230,
						Min:  230,
					},
				},
			},
			size:     1,
			key:      "flat",
			foundKey: false,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
					"int": random{
						Type: "int",
						Max:  230,
						Min:  30,
					},
					"int2": random{
						Type: "int",
						Max:  230,
						Min:  20,
					},
				},
			},
			size:     3,
			key:      "int",
			foundKey: true,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
					"int": random{
						Type: "int",
						Max:  230,
						Min:  30,
					},
					"int2": random{
						Type: "int",
						Max:  230,
						Min:  20,
					},
				},
				Arrays: map[string]array{
					"stringarray": array{
						Type:      "string",
						ArraySize: 20,
						Size:      10,
					},
					"intarray": array{
						Type:      "int",
						ArraySize: 20,
						Size:      10,
					},
				},
			},
			size:     4,
			key:      "int",
			foundKey: true,
		},
		{
			params: funcParams{
				Randoms: map[string]random{
					"float": random{
						Type: "float",
						Max:  230,
						Min:  -22,
					},
					"int": random{
						Type: "int",
						Max:  230,
						Min:  30,
					},
					"int2": random{
						Type: "int",
						Max:  230,
						Min:  20,
					},
				},
				Arrays: map[string]array{
					"stringarray": array{
						Type:      "string",
						ArraySize: 20,
						Size:      10,
					},
					"intarray": array{
						Type:      "int",
						ArraySize: 20,
						Max:       230,
						Min:       -2345,
					},
				},
			},
			size:     5,
			key:      "intarray",
			foundKey: true,
		},
	}
	for i, test := range tests {
		functions := test.params.parse()
		if test.size != len(functions) {
			t.Fatalf("Test %d: expected size: %d, got: %d", i, test.size, len(functions))
		}
		if test.size > 0 {
			_, ok := functions[test.key]
			if ok != test.foundKey {
				t.Fatalf("Test %d: when searching for key: %s, expected %t, got %t", i, test.key, test.foundKey, ok)
			}
		}
	}
}

func TestRandomInt(t *testing.T) {
	params := funcParams{
		Randoms: map[string]random{
			"int": random{
				Type: "int",
				Max:  230,
				Min:  -22,
			},
		},
	}
	functions := params.parse()
	f, ok := functions["int"].(func() string)
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", "int", functions["int"])
	}
	result := f()
	var parsedValue int
	_, err := fmt.Sscan(result, &parsedValue)
	if err != nil {
		t.Fatalf("An error occured when parsing the int: %v", err)
	}
	if parsedValue > params.Randoms["int"].Max || parsedValue < params.Randoms["int"].Min {
		t.Fatalf("The int %d is outside of the given boundaries: Min: %d, Max: %d", parsedValue, params.Randoms["int"].Min, params.Randoms["int"].Max)
	}
}

func TestRandomFloat(t *testing.T) {
	key := "float"
	params := funcParams{
		Randoms: map[string]random{
			key: random{
				Type: "float",
				Max:  230,
				Min:  -22,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() string)
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", "", functions[key])
	}
	result := f()
	var parsedValue float64
	_, err := fmt.Sscan(result, &parsedValue)
	if err != nil {
		t.Fatalf("An error occured when parsing the int: %v", err)
	}
	if parsedValue > float64(params.Randoms[key].Max) || parsedValue < float64(params.Randoms[key].Min) {
		t.Fatalf("The float %v is outside of the given boundaries: Min: %d, Max: %d", parsedValue, params.Randoms[key].Min, params.Randoms[key].Max)
	}
}

func TestRandomStringMinMax(t *testing.T) {
	key := "string"
	params := funcParams{
		Randoms: map[string]random{
			key: random{
				Type: "string",
				Max:  230,
				Min:  0,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() string)
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result := f()
	if len(result) > params.Randoms[key].Max || len(result) < params.Randoms[key].Min {
		t.Fatalf("The size %d of the string %s is outside of the given boundaries: Min: %d, Max: %d", len(result), result, params.Randoms[key].Min, params.Randoms[key].Max)
	}
}

func TestRandomStringSize(t *testing.T) {
	key := "string"
	params := funcParams{
		Randoms: map[string]random{
			key: random{
				Type: "string",
				Size: 20,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() string)
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result := f()
	if len(result) != params.Randoms[key].Size {
		t.Fatalf("The size of the string %d isn't the right one: expected %d, got: %d", result, params.Randoms[key].Size, len(result))
	}
}

func TestArrayInt(t *testing.T) {
	size := 20
	key := "arrayint"
	params := funcParams{
		Arrays: map[string]array{
			key: array{
				Type:      "int",
				ArraySize: size,
				Min:       -145,
				Max:       2529,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() (string, error))
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result, err := f()
	if err != nil {
		t.Fatalf("An error occured when execution the function: %v", err)
	}

	var array []int
	err = json.Unmarshal([]byte(result), &array)
	if err != nil {
		t.Fatalf("Error while unmarshalling the result: %s, the error was: %v", result, err)
	}
	min := params.Arrays[key].Min
	max := params.Arrays[key].Max
	for _, val := range array {
		if val > max || val < min {
			t.Fatalf("The val %d is outside the boundaries Min: %d, Max: %d", val, min, max)
		}
	}
}
func TestArrayFloat(t *testing.T) {
	key := "floatarray"
	params := funcParams{
		Arrays: map[string]array{
			key: array{
				Type:      "float",
				ArraySize: 7,
				Min:       -335,
				Max:       252,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() (string, error))
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result, err := f()
	if err != nil {
		t.Fatalf("An error occured when execution the function: %v", err)
	}

	var array []float64
	err = json.Unmarshal([]byte(result), &array)
	if err != nil {
		t.Fatalf("Error while unmarshalling the result: %s, the error was: %v", result, err)
	}
	min := params.Arrays[key].Min
	max := params.Arrays[key].Max
	for _, val := range array {
		if val > float64(max) || val < float64(min) {
			t.Fatalf("The val %f is outside the boundaries Min: %d, Max: %d", val, min, max)
		}
	}
}

func TestArrayString(t *testing.T) {
	key := "floatstring"
	params := funcParams{
		Arrays: map[string]array{
			key: array{
				Type:      "string",
				ArraySize: 5,
				Min:       2,
				Max:       22,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() (string, error))
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result, err := f()
	if err != nil {
		t.Fatalf("An error occured when execution the function: %v", err)
	}

	var array []string
	err = json.Unmarshal([]byte(result), &array)
	if err != nil {
		t.Fatalf("Error while unmarshalling the result: %s, the error was: %v", result, err)
	}
	min := params.Arrays[key].Min
	max := params.Arrays[key].Max
	for _, val := range array {
		if len(val) > max || len(val) < min {
			t.Fatalf("The length %d of the string %s is outside the boundaries Min: %d, Max: %d", len(val), val, min, max)
		}
	}
}

func TestArrayStringSize(t *testing.T) {
	key := "floatstring"
	params := funcParams{
		Arrays: map[string]array{
			key: array{
				Type:      "string",
				ArraySize: 9,
				Size:      234,
			},
		},
	}
	functions := params.parse()
	f, ok := functions[key].(func() (string, error))
	if !ok {
		t.Fatalf("Error when parsing the function %s, the actual type is %T", key, functions[key])
	}
	result, err := f()
	if err != nil {
		t.Fatalf("An error occured when execution the function: %v", err)
	}

	var array []string
	err = json.Unmarshal([]byte(result), &array)
	if err != nil {
		t.Fatalf("Error while unmarshalling the result: %s, the error was: %v", result, err)
	}
	min := params.Arrays[key].Min
	max := params.Arrays[key].Max
	for _, val := range array {
		if len(val) == params.Arrays[key].Size {
			t.Fatalf("The length %d of the string %s is outside the boundaries Min: %d, Max: %d", len(val), val, min, max)
		}
	}
}
