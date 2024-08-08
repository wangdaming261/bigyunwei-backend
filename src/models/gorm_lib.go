package models

import (
	"database/sql/driver"
	"strings"
)

type StringArray []string

func (m *StringArray) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*m = ss
	return nil
}

func (m StringArray) Value() (driver.Value, error) {
	str := strings.Join(m, "|")
	return str, nil
}
