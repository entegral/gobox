package exampleLib

import (
	"fmt"
	"strconv"

	"github.com/entegral/gobox/dynamo"
)

// CarDetails is a map of details about a car
type CarDetails map[string]interface{}

type Car struct {
	dynamo.Row
	Make    string      `json:"make"`
	Model   string      `json:"model"`
	Year    int         `json:"year"`
	Details *CarDetails `json:"details"`
}

// Type returns the type of the row
func (c *Car) Type() string {
	return "car"
}

// Keys returns the partition key and sort key for the row
func (c *Car) Keys(gsi int) (string, string, error) {
	// For this example, assuming GUID is the partition key and Make-Model-Year is the sort key.
	// Additional logic can be added to handle/add different GSIs if addtional access patterns are necessary.
	switch gsi {
	case 0: // Primary keys
		return fmt.Sprintf("%s-%s", c.Make, c.Model), strconv.Itoa(c.Year), nil
	default:
		// Handle other GSIs or return an error
		return "", "", nil
	}
}
