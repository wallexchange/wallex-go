package wallex

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Number represents a number or a number string.
type Number string

// IsUndefined returns true if n is undefined or not a valid number.
func (n Number) IsUndefined() bool {
	return n == ""
}

// Float converts n to a float value.
func (n Number) Float() float64 {
	f, _ := strconv.ParseFloat(string(n), 64)
	return f
}

// UnmarshalJSON deserializes n from JSON data.
func (n *Number) UnmarshalJSON(data []byte) error {
	*n = ""
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			*n = Number(s)
		}
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*n = Number(fmt.Sprintf("%f", f))
		return nil
	}
	return nil
}
