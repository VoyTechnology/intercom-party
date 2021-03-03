package distance

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := map[string]struct {
		want int
		err  error
	}{
		"1m":       {1, nil},
		"10m":      {10, nil},
		"100m":     {100, nil},
		"1000m":    {1000, nil},
		"1km":      {1000, nil},
		"10km":     {10000, nil},
		"100km":    {100000, nil},
		"1km1m":    {1001, nil},
		"1km100m":  {1100, nil},
		"1km1000m": {2000, nil},
		"invalid":  {0, ErrInvalidDistance},
		"1":        {0, ErrInvalidDistance},
		"m":        {0, ErrInvalidDistance},
	}

	for in, tc := range testCases {
		t.Run(in, func(t *testing.T) {
			got, err := Parse(in)
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("Parse(%s) err = %v, want %v", in, err, tc.err)
				}
				return
			}

			if got != tc.want {
				t.Errorf("Parse(%s) = %v, want %v", in, got, tc.want)
			}
		})
	}
}

func TestDistance(t *testing.T) {
	testCases := map[float64][]float64{
		0:         {0, 0, 0, 0},
		45711.259: {52.986375, -6.043701, 53.339428, -6.257664}, // Christina McArdle
	}

	for want, in := range testCases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			got := Distance(in[0], in[1], in[2], in[3])
			if !equal(got, want) {
				t.Errorf("Distance() = %v, want %v", got, want)
			}
		})
	}
}

func TestDtor(t *testing.T) {
	testCases := map[float64]float64{
		0:   0,
		180: math.Pi,
		360: 2 * math.Pi,
		90:  math.Pi / 2,
		45:  math.Pi / 4,
	}

	for in, want := range testCases {
		t.Run(fmt.Sprint(in), func(t *testing.T) {
			if got := dtor(in); got != want {
				t.Errorf("dtor(%f) = %f, want %f", in, got, want)
			}
		})
	}
}

// Float equality check within a threshold. We are dealing with meters here, so
// we don't care about extreme precision.
func equal(a, b float64) bool {
	return math.Abs(a-b) <= 0.01
}
