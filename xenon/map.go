package xenon

import (
	"encoding/json"
	"strconv"
)

type Map map[string]interface{}

func ToMap(bytes []byte) (Map, error) {
	var data Map
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MustMap(bytes []byte) Map {
	data, err := ToMap(bytes)
	PanicNotNilError(err)
	return data
}

func ToMaps(bytes []byte) ([]Map, error) {
	var data []Map
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MustMaps(bytes []byte) []Map {
	data, err := ToMaps(bytes)
	PanicNotNilError(err)
	return data
}

func (this Map) Clone() Map {
	return MustMap(this.MustJSON())
}

func (this Map) ToJSON() ([]byte, error) {
	return json.Marshal(this)
}

func (this Map) MustJSON() []byte {
	data, err := this.ToJSON()
	PanicNotNilError(err)
	return data
}

func (this Map) Set(key string, val interface{}) Map {
	this[key] = val
	return this
}

func (this Map) Del(key string) Map {
	delete(this, key)
	return this
}

func (this Map) DeepMerge(src Map) Map {
	for key, srcVal := range src {
		if distVal, exists := this[key]; exists {
			if dstMap, ok1 := distVal.(map[string]interface{}); ok1 {
				if srcMap, ok2 := srcVal.(map[string]interface{}); ok2 {
					Map(dstMap).DeepMerge(srcMap)
					continue
				}
				if srcMap, ok2 := srcVal.(Map); ok2 {
					Map(dstMap).DeepMerge(srcMap)
					continue
				}
			}

			if dstMap, ok1 := distVal.(Map); ok1 {
				if srcMap, ok2 := srcVal.(Map); ok2 {
					dstMap.DeepMerge(srcMap)
					continue
				}
				if srcMap, ok2 := srcVal.(map[string]interface{}); ok2 {
					dstMap.DeepMerge(srcMap)
					continue
				}
			}
		}
		this[key] = srcVal
	}
	return this
}

func (this Map) GetBool(key string, val ...bool) bool {
	switch value := this[key].(type) {
	case bool:
		return value
	case string:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return false
		}
		return v
	}
	if len(val) > 0 {
		return val[0]
	}
	return false
}

func (this Map) GetString(key string, val ...string) string {
	if data, ok := this[key].(string); ok {
		return data
	}
	if len(val) > 0 {
		return val[0]
	}
	return ""
}

func (this Map) GetInt(key string, val ...int) int {
	switch value := this[key].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case int32:
		return int(value)
	case float64:
		return int(value)
	case float32:
		return int(value)
	case string:
		v, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}
		return v
	}
	if len(val) > 0 {
		return val[0]
	}
	return 0
}

func (this Map) GetFloat(key string, val ...float64) float64 {
	switch value := this[key].(type) {
	case int:
		return float64(value)
	case int64:
		return float64(value)
	case int32:
		return float64(value)
	case float64:
		return float64(value)
	case float32:
		return float64(value)
	case string:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}
		return v
	}
	if len(val) > 0 {
		return val[0]
	}
	return 0
}

func (this Map) GetMap(key string, val ...Map) Map {
	v := this[key]
	switch v.(type) {
	case Map:
		return v.(Map)
	case map[string]interface{}:
		return v.(map[string]interface{})
	}
	if len(val) > 0 {
		return val[0]
	}
	return nil
}

func (this Map) GetMaps(key string, val ...[]Map) []Map {
	v := this[key]
	switch v.(type) {
	case []Map:
		return v.([]Map)
	case []interface{}:
		maps := make([]Map, 0)
		for _, value := range v.([]interface{}) {
			switch value.(type) {
			case Map:
				maps = append(maps, value.(Map))
			case map[string]interface{}:
				maps = append(maps, Map(value.(map[string]interface{})))
			}
		}
		return maps
	}
	if len(val) > 0 {
		return val[0]
	}
	return nil
}

func (this Map) GetSlice(key string, val ...[]interface{}) []interface{} {
	v := this[key]
	switch v.(type) {
	case []Map:
		interfaces := make([]interface{}, 0)
		for _, value := range v.([]Map) {
			interfaces = append(interfaces, value)
		}
		return interfaces
	case []interface{}:
		return v.([]interface{})
	}
	if len(val) > 0 {
		return val[0]
	}
	return nil
}