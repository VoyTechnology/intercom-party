package main // import "github.com/voytechnology/intercom-party/cmd/party"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/voytechnology/intercom-party/internal/customer"
	"github.com/voytechnology/intercom-party/internal/distance"
	"github.com/voytechnology/intercom-party/internal/office"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

var (
	officeFlag   = flag.String("office", "Dublin", "Office to get the party from")
	distanceFlag = flag.String("distance", "100km", "maximum distance for the customers to be invited")
)

func run() error {
	flag.Parse()

	return customersInOfficeRadius(
		os.Stdin, os.Stdout, *officeFlag, *distanceFlag, distance.Distance)
}

func customersInOfficeRadius(r io.Reader, w io.Writer, officeName, dist string, distFn distanceFunc) error {
	officeCoords, err := office.Coordinates(officeName)
	if err != nil {
		return fmt.Errorf("unable to get office coordinates for office %s: %w",
			officeName, err)
	}

	maxDistance, err := distance.Parse(dist)
	if err != nil {
		return fmt.Errorf("unable to parse the distance from office: %w", err)
	}

	customers, err := customer.Parse(r)
	if err != nil {
		return fmt.Errorf("unable to parse customers: %w", err)
	}

	invited := customer.Filter(customers, withinRadius(distFn, officeCoords, maxDistance))

	// sort all customers in ascending order by their customer ID.
	sort.Sort(customer.ByID(invited))

	if _, err := customer.Write(w, invited); err != nil {
		return fmt.Errorf("unable to write customers: %w", err)
	}

	return nil
}

func withinRadius(distFn distanceFunc, from []float64, max int) func(c customer.Customer) bool {
	return func(c customer.Customer) bool {
		return distFn(
			c.Latitude,
			c.Longitude,
			from[0],
			from[1]) < float64(max)
	}
}

type distanceFunc func(x1, y1, x2, y2 float64) float64
