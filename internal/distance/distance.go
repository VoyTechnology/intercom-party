// Package distance provides utilities for handing distances.
package distance

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const (
	// Average radius of the earth in meters
	EarthRadius = 6371009
)

var (
	ErrInvalidDistance = errors.New("invalid distance")
)

// Parse converts from human readable strings to meters. The acceptable
// format is 1km, 1km100m, 100m, where the numbers can be changed. The
// capitalization doesn't matter. If there is an error in parsing an error is
// returned.
func Parse(s string) (int, error) {
	// first we just convert it to lowercase
	s = strings.ToLower(s)

	// If the len of the string is less than 2, its definately invalid
	if len(s) < 2 {
		return 0, ErrInvalidDistance
	}

	distance := 0

	var n string
	for i := 0; i < len(s); i++ {
		m := 1
		switch {
		case unicode.IsDigit(rune(s[i])):
			n += string(s[i])
			continue
		case s[i] == 'k':
			m = 1000
			i++ // skip m
		case s[i] == 'm':
			m = 1
		}

		v, err := strconv.Atoi(n)
		if err != nil {
			return 0, ErrInvalidDistance
		}
		distance += v * m
		n = ""
	}

	return distance, nil
}

// Distance calculates the distance between two points on a globe.
// Adapted from https://en.wikipedia.org/wiki/Great-circle_distance
// Because we are operating on float64, we do not care about inprecision errors
// as we care about very small distances.
func Distance(x1, y1, x2, y2 float64) float64 {
	// convert to radians
	x1, y1, x2, y2 = dtor(x1), dtor(y1), dtor(x2), dtor(y2)

	lonDiff := math.Abs(x2 - x1)

	a := math.Sin(y1) * math.Sin(y2)
	b := math.Cos(y1) * math.Cos(y2) * math.Cos(lonDiff)
	rd := math.Acos(a + b)

	return rd * EarthRadius
}

// dtor converts degrees to radians
func dtor(d float64) float64 {
	return (d * math.Pi) / 180
}
