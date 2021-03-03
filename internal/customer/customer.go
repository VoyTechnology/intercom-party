// Package customer contains the definition for the customer, as well as the
// parsing of each customers from reader
package customer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Customer contains the fields for the customer
type Customer struct {
	Latitude  float64
	Longitude float64
	UserID    int
	Name      string
}

var _ json.Unmarshaler = (*Customer)(nil)

// UnmarshalJSON helps to parse the customer data, as they are formated in
// non-standard format like floats being strings.
func (c *Customer) UnmarshalJSON(b []byte) (err error) {
	var invalid customer
	if err := json.Unmarshal(b, &invalid); err != nil {
		return err
	}

	// Handle the easy cases first
	c.Name = invalid.Name
	c.UserID = invalid.UserID

	c.Latitude, err = strconv.ParseFloat(invalid.Latitude, 64)
	if err != nil {
		return fmt.Errorf("unable to convert lattitude: %w", err)
	}

	c.Longitude, err = strconv.ParseFloat(invalid.Longitude, 64)
	if err != nil {
		return fmt.Errorf("unable to convert longitude: %w", err)
	}

	return nil
}

var _ json.Marshaler = (*Customer)(nil)

// MarshalJSON converts the customers back to the input specific format.
func (c Customer) MarshalJSON() ([]byte, error) {
	return json.Marshal(customer{
		UserID:    c.UserID,
		Name:      c.Name,
		Longitude: fmt.Sprint(c.Longitude),
		Latitude:  fmt.Sprint(c.Latitude),
	})
}

// customer contains the non-standard encoding of the customer data
type customer struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// Parse the customers from a reader
func Parse(r io.Reader) ([]Customer, error) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	var customers []Customer

	for s.Scan() {
		var customer Customer
		if err := json.Unmarshal(s.Bytes(), &customer); err != nil {
			return nil, fmt.Errorf("unable to parse customer: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// Filter only returns the customers which match the filter function f.
// TODO: This re-allocates for all customers. maybe just remove from the
//       customer list in the future, and return the same slice.
func Filter(cs []Customer, f func(c Customer) bool) []Customer {
	out := make([]Customer, 0, len(cs))
	for _, c := range cs {
		if f(c) {
			out = append(out, c)
		}
	}

	return out
}

// Write all customers to a io.Writer in the same format as the input.
func Write(w io.Writer, cs []Customer) (int, error) {
	total := 0
	for _, c := range cs {
		b, err := json.Marshal(c)
		if err != nil {
			return total, err
		}

		n, err := fmt.Fprintln(w, string(b))
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// ByID allows to sort the customer data by their ID, in ascending order.
type ByID []Customer

func (c ByID) Len() int           { return len(c) }
func (c ByID) Less(i, j int) bool { return c[i].UserID < c[j].UserID }
func (c ByID) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
