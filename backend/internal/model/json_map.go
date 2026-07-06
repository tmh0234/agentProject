package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]string

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = JSONMap{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported JSONMap type %T", value)
	}
	if len(bytes) == 0 {
		*m = JSONMap{}
		return nil
	}
	return json.Unmarshal(bytes, m)
}
