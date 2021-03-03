package customer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestUnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		want Customer
		err  error
	}{
		`{"latitude": "1.1", "user_id": 1, "name": "Test Case", "longitude": "1.1"}`: {
			Customer{Latitude: 1.1, UserID: 1, Name: "Test Case", Longitude: 1.1},
			nil,
		},
		`{"latitude": "somewhere", "user_id": 1, "name": "Test Case", "longitude": "1.1"}`: {
			Customer{Latitude: 1.1, UserID: 1, Name: "Test Case", Longitude: 1.1},
			strconv.ErrSyntax,
		},
		`{"latitude": "1.1", "user_id": 1, "name": "Test Case", "longitude": "somewhere"}`: {
			Customer{Latitude: 1.1, UserID: 1, Name: "Test Case", Longitude: 1.1},
			strconv.ErrSyntax,
		},
		// JSON package doesn't work well with errors.Is. In normal
		// circumstances this error would have worked perfectly. Or if we
		// do a hacky test and compare Error strings themselves. This test is
		// disabled because of this, but looking at the error message by hand
		// shows exactly the same result:
		//     Unmarshal() err = json: cannot unmarshal number into Go struct field customer.latitude of type string, want json: cannot unmarshal number into Go struct field customer.latitude of type string
		//
		// `{"latitude": 1.1, "user_id": 1, "name": "Test Case", "longitude": "1.1"}`: {
		// 	Customer{Latitude: 1.1, UserID: 1, Name: "Test Case", Longitude: 1.1},
		// 	&json.UnmarshalTypeError{
		// 		Value:  "number",
		// 		Struct: "customer",
		// 		Field:  "latitude",
		// 		Offset: 16,
		// 		Type:   reflect.TypeOf(""),
		// 	},
		// },
	}

	for in, tc := range testCases {
		t.Run(in, func(t *testing.T) {
			var got Customer
			if err := json.Unmarshal([]byte(in), &got); err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("Unmarshal() err = %v, want %v", err, tc.err)
				}
				return
			}

			if got != tc.want {
				t.Errorf("Unmarshal() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	testCases := map[*Customer]string{
		{Name: "Test", UserID: 1, Longitude: 1.1, Latitude: 1.1}: `{"user_id":1,"name":"Test","longitude":"1.1","latitude":"1.1"}`,
	}

	for in, want := range testCases {
		t.Run(fmt.Sprint(in), func(t *testing.T) {
			res, err := json.Marshal(in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(res); got != want {
				t.Errorf("json.Marshal() = %s, want %s", got, want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want []Customer
		err  error
	}{
		"empty": {``, nil, nil},
		"single": {
			`{"latitude": "1.1", "user_id": 1, "name": "Test Case", "longitude": "1.1"}`,
			[]Customer{
				{Name: "Test Case", Latitude: 1.1, Longitude: 1.1, UserID: 1},
			},
			nil,
		},
		"multiple": {
			`{"latitude": "1.1", "user_id": 1, "name": "Test Case", "longitude": "1.1"}
			{"latitude": "1.1", "user_id": 2, "name": "Case Test", "longitude": "1.1"}`,
			[]Customer{
				{Name: "Test Case", Latitude: 1.1, Longitude: 1.1, UserID: 1},
				{Name: "Case Test", Latitude: 1.1, Longitude: 1.1, UserID: 2},
			},
			nil,
		},
		"error": {
			`{"latitude": "invalid", "user_id": 1, "name": "Test Case", "longitude": "1.1"}`,
			nil,
			strconv.ErrSyntax,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := strings.NewReader(tc.in)
			got, err := Parse(r)
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("Parse() err = %v, want %v", err, tc.err)
				}
				return
			}

			if diff := deep.Equal(got, tc.want); diff != nil {
				t.Errorf("Parse() = %v, want %v, diff = %v", got, tc.want, diff)
			}
		})

	}
}

func TestFilter(t *testing.T) {
	testCases := map[string]struct {
		in   []Customer
		want []Customer
	}{
		"empty": {[]Customer{}, []Customer{}},
		"remove": {
			[]Customer{{UserID: 1}},
			[]Customer{},
		},
		"keep": {
			[]Customer{{UserID: 2}},
			[]Customer{{UserID: 2}},
		},
		"mix": {
			[]Customer{{UserID: 1}, {UserID: 2}},
			[]Customer{{UserID: 2}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Only let customers with IDs of even numbers
			f := func(c Customer) bool {
				return c.UserID%2 == 0
			}
			got := Filter(tc.in, f)

			if diff := deep.Equal(got, tc.want); diff != nil {
				t.Errorf("Filter() = %v, want %v, diff = %v",
					got, tc.want, diff)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	testCases := map[string]struct {
		rw   io.ReadWriter
		cs   []Customer
		want string
		err  error
	}{
		"empty": {
			new(bytes.Buffer),
			nil,
			"",
			nil,
		},
		"good": {
			new(bytes.Buffer),
			[]Customer{{UserID: 1, Name: "Test"}},
			`{"user_id":1,"name":"Test","longitude":"0","latitude":"0"}`,
			nil,
		},
		"bad writer": {
			&badReadWriter{},
			[]Customer{{UserID: 1, Name: "Test"}},
			"",
			errBadWrite,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if _, err := Write(tc.rw, tc.cs); err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("Write() err = %v, want %v", err, tc.err)
				}
				return
			}

			res, err := io.ReadAll(tc.rw)
			if err != nil {
				t.Fatalf("unexpected error when reading back from writer: %v",
					err)
			}
			if got := strings.TrimSpace(string(res)); got != tc.want {
				t.Errorf("Write() = %s, want %s", got, tc.want)
			}
		})

	}

}

func TestSortByID(t *testing.T) {
	testCases := map[string]struct {
		in   []Customer
		want []Customer
	}{
		"empty": {},
		"single": {
			[]Customer{{UserID: 3}},
			[]Customer{{UserID: 3}},
		},
		"multiple": {
			[]Customer{
				{UserID: 3},
				{UserID: 1},
				{UserID: 2},
			},
			[]Customer{
				{UserID: 1},
				{UserID: 2},
				{UserID: 3},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sort.Sort(ByID(tc.in))

			if diff := deep.Equal(tc.in, tc.want); diff != nil {
				t.Errorf("sort(ByID) diff = %v", diff)
			}
		})

	}
}

var errBadWrite = errors.New("bad write")

type badReadWriter struct{}

func (badReadWriter) Write(_ []byte) (int, error) {
	return 0, errBadWrite
}

func (badReadWriter) Read(_ []byte) (int, error) {
	return 0, nil
}
