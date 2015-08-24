package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/evermax/iseva/util"
)

type funcParams struct {
	Randoms map[string]random `json:"rand"`
	Arrays  map[string]array  `json:"array"`
}

type random struct {
	Type string `json:"type"`
	Size int    `json:"size"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

type array struct {
	Type      string `json:"type"`
	ArraySize int    `json:"arraysize"`
	Size      int    `json:"size"`
	Min       int    `json:"min"`
	Max       int    `json:"max"`
}

func (fp funcParams) parse() (fcts map[string]interface{}) {
	fcts = make(map[string]interface{})
	for name, randomParam := range fp.Randoms {
		if randomParam.Max > randomParam.Min {
			switch randomParam.Type {
			case "string":
				if randomParam.Min >= 0 {
					length := randomParam.Min + rand.Intn(randomParam.Max-randomParam.Min)
					fcts[name] = func() string {
						return util.RandString(length)
					}
				}
			case "int":
				val := randomParam.Min + rand.Intn(randomParam.Max-randomParam.Min)
				fcts[name] = func() string {
					return fmt.Sprintf("%d", val)
				}
			case "float":
				fcts[name] = func() string {
					return fmt.Sprintf("%f", float64(randomParam.Min)+rand.Float64()*float64(randomParam.Max-randomParam.Min))
				}
			}
		} else if randomParam.Size > 0 {
			if randomParam.Type == "string" {
				fcts[name] = func() string {
					return util.RandString(randomParam.Size)
				}
			}
		}
	}
	for name, arrParam := range fp.Arrays {
		if arrParam.ArraySize > 0 {
			if arrParam.Max > arrParam.Min {
				switch arrParam.Type {
				case "string":
					size := arrParam.ArraySize
					min := arrParam.Min
					max := arrParam.Max
					fcts[name] = func() (string, error) {
						var array = make([]string, size)
						for i := 0; i < size; i++ {
							length := arrParam.Min + rand.Intn(max-min)
							array[i] = util.RandString(length)
						}
						arr, err := json.Marshal(array)
						if err != nil {
							return "", err
						}
						return string(arr), nil
					}
				case "int":
					size := arrParam.ArraySize
					min := arrParam.Min
					max := arrParam.Max
					fcts[name] = func() (string, error) {
						var array = make([]int, size)
						for i := 0; i < size; i++ {
							array[i] = min + rand.Intn(max-min)
						}
						arr, err := json.Marshal(array)
						if err != nil {
							return "", err
						}
						return string(arr), nil
					}
				case "float":
					size := arrParam.ArraySize
					min := arrParam.Min
					max := arrParam.Max
					fcts[name] = func() (string, error) {
						var array = make([]float64, size)
						for i := 0; i < size; i++ {
							array[i] = float64(min) + rand.Float64()*float64(max-min)
						}
						arr, err := json.Marshal(array)
						if err != nil {
							return "", err
						}
						return string(arr), nil
					}
				}
			} else if arrParam.Size > 0 {
				if arrParam.Type == "string" {
					size := arrParam.ArraySize
					fcts[name] = func() (string, error) {
						var array = make([]string, size)
						for i := 0; i < size; i++ {
							array[i] = util.RandString(size)
						}
						arr, err := json.Marshal(array)
						if err != nil {
							return "", err
						}
						return string(arr), nil
					}
				}
			}
		}
	}
	return
}
