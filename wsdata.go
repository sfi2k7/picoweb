package picoweb

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Interface interface{}
type WsData map[string]Interface

type Subscribers map[string]struct{}

func (wd WsData) Get(k string) Interface {
	v := wd[k]
	return v
}

func (wd WsData) Set(k string, v Interface) {
	wd[k] = v
}

func (wd WsData) Json() string {
	bts, err := json.Marshal(wd)
	if err != nil {
		fmt.Println("Error in WSDATA JSON()", err)
	}
	return string(bts)
}

func WsDataFromString(v string) *WsData {
	var single WsData
	json.Unmarshal([]byte(v), &single)
	return &single
}

func WsDataFromMapString(m map[string]string) WsData {
	var data = WsData{}
	if m == nil {
		return data
	}

	for k, v := range m {
		data[k] = v
	}
	return data
}

func (wd WsData) Remove(k string) {
	delete(wd, k)
}

func (wd WsData) Bool(k string) bool {
	v := wd.Get(k)

	if b, ok := v.(bool); ok {
		// fmt.Println("Returning bool", b)
		return b
	}

	if str, ok := v.(string); ok {
		if len(str) == 0 {
			return false
		}
		// fmt.Println("Parsing as bool", str)
		b, err := strconv.ParseBool(str)
		if err != nil {
			fmt.Println("Error parding bool")
		}
		return b
	}
	// fmt.Println("Returning false (default)")
	return false
}

func (wd WsData) Clone() WsData {
	cloned := WsData{}
	for k, v := range wd {
		cloned[k] = v
	}
	return cloned
}

func (wd WsData) DataAsString(k string) string {
	data, ok := wd[k]
	if !ok {
		return ""
	}

	switch tp := data.(type) {
	case nil:
		return ""
	case map[string]interface{}:
		b, _ := json.Marshal(tp)
		return string(b)
	case string:
		return tp
	default:
		return fmt.Sprint(tp)
	}
}

func (wd WsData) String(k string) string {
	v := wd[k]
	if v == nil {
		return ""
	}

	str, ok := v.(string)
	if ok {
		return str
	}
	return fmt.Sprint(v)
}

func (wd WsData) ArrayString(k string) []string {
	var result []string
	switch tp := wd[k].(type) {
	case []string:
		result = tp
	case []interface{}:
		for _, i := range tp {
			s, ok := i.(string)
			if ok {
				result = append(result, s)
			}
		}
	case []int:
	case []int64:
	case []float64:
		for _, i := range tp {
			result = append(result, fmt.Sprint(i))
		}
	}

	return result
}

func (wd WsData) Int(k string) int {
	// fmt.Println("WD", wd)
	v := wd[k]
	// fmt.Println("WD V", v)
	if v == nil {
		fmt.Println("Int:V us nil")
		return 0
	}

	switch vt := v.(type) {
	case int:
		return vt
	case int64:
		return int(vt)
	case float64:
		return int(vt)
	default:
		return -1
	}
}
