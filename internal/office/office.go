// Package offices contains the list of offices and their location.
package office

import (
	"errors"
	"strings"
)

var (
	ErrInvalidOffice = errors.New("office does not (yet) exist")
)

// offices lists contains all the  offices and their location. Because the list
// will rarely change, we can just edit it in code rather than dynamically load
// it. If the frequency of how often offices open changes, perhaps it would be
// better to load it dynamically from some office inventory system?
var offices = map[string][]float64{
	"Dublin": {53.339428, -6.257664},
}

// Coordinates returns the coordinates for the office if it exists, or an error if
// the office does not yet exist.
func Coordinates(name string) ([]float64, error) {
	if coords, exists := offices[strings.Title(name)]; exists {
		return coords, nil
	} else {
		return nil, ErrInvalidOffice
	}
}
