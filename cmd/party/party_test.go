package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/voytechnology/intercom-party/internal/customer"
	"github.com/voytechnology/intercom-party/internal/distance"
	"github.com/voytechnology/intercom-party/internal/office"
)

func TestCustomersInOfficeRadius_Within(t *testing.T) {
	// ~41.75km from the office
	want := `{"user_id":12,"name":"Christina McArdle","longitude":"-6.043701","latitude":"52.986375"}`
	in := strings.NewReader(want)

	out := new(bytes.Buffer)

	if err := customersInOfficeRadius(in, out, "Dublin", "100km"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want += "\n"

	if got := out.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCustomersInOfficeRadius_Outside(t *testing.T) {
	// ~41.75km from the office
	in := strings.NewReader(`{"latitude": "52.986375", "user_id": 12, "name": "Christina McArdle", "longitude": "-6.043701"}`)

	out := new(bytes.Buffer)

	if err := customersInOfficeRadius(in, out, "Dublin", "30km"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ""

	if got := out.String(); got != want {
		t.Errorf("got %v, want %v", got, "")
	}
}

func TestCustomersInOfficeRadius_BadOffice(t *testing.T) {
	in, out := new(bytes.Buffer), new(bytes.Buffer)

	if err := customersInOfficeRadius(in, out, "Bad Office", "100km"); !errors.Is(err, office.ErrInvalidOffice) {
		t.Fatalf("bad office err = %v, want %v", err, office.ErrInvalidOffice)
	}
}

func TestCustomersInOfficeRadius_BadDistance(t *testing.T) {
	in, out := new(bytes.Buffer), new(bytes.Buffer)

	if err := customersInOfficeRadius(in, out, "Dublin", "bad distance"); !errors.Is(err, distance.ErrInvalidDistance) {
		t.Fatalf("bad office err = %v, want %v", err, distance.ErrInvalidDistance)
	}
}

func TestCustomersInOfficeRadius_BadCustomers(t *testing.T) {
	// ~41.75km from the office
	in := strings.NewReader(`{"latitude": "1.1", "user_id": 12, "name": "Christina McArdle", "longitude": "-6.043701"}`)

	out := &badWriter{}

	if err := customersInOfficeRadius(in, out, "Dublin", "30km"); !errors.Is(err, errBadWriter) {
		t.Fatalf("unexpected error: %v, want %v", err, strconv.ErrSyntax)
	}
}

func TestWithinRadius(t *testing.T) {
	testCases := map[int]bool{
		1:   true,
		5:   true,
		11:  false,
		100: false,
	}

	for in, want := range testCases {
		t.Run(fmt.Sprint(in), func(t *testing.T) {
			distFn := func(a, b, c, d float64) float64 { return float64(in) }
			if got := withinRadius(distFn, []float64{0, 0}, 10)(customer.Customer{}); got != want {
				t.Errorf("withingRadius() = %t, want %t", got, want)
			}
		})
	}
}

var errBadWriter = errors.New("bad writer write error")

type badWriter struct{}

func (*badWriter) Write(_ []byte) (int, error) {
	return 0, errBadWriter
}
