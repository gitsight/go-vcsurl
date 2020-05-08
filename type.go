package vcsurl

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements database/sql.Scanner.
func (x *Kind) Scan(v interface{}) error {
	if data, ok := v.([]byte); ok {
		*x = Kind(data)
		return nil
	}
	return fmt.Errorf("%T.Scan failed: %v", x, v)
}

// Scan implements database/sql/driver.Valuer
func (x Kind) Value() (driver.Value, error) {
	return string(x), nil
}
