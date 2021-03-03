package office

import (
	"errors"
	"testing"
)

func TestCoordinates(t *testing.T) {
	testCases := map[string]error{
		"Dublin": nil,
		"Mars":   ErrInvalidOffice,
	}

	for in, wantErr := range testCases {
		t.Run(in, func(t *testing.T) {
			// We don't actually care about the coordinates, just if we
			// are matching the cities.
			if _, err := Coordinates(in); err != nil {
				if !errors.Is(err, wantErr) {
					t.Errorf("Coordinates() err = %v, want %v", err, wantErr)
				}
				return
			}
		})

	}
}
