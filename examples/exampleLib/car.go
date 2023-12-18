package exampleLib

import (
	"fmt"

	"github.com/entegral/gobox/dynamo"
)

type Car struct {
	dynamo.Row
	OwnerEmail string `json:"ownerEmail"`
	Make       string `json:"make"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
}

// Type returns the type of the row
func (c *Car) Type() string {
	return "car"
}

// Keys returns the partition key and sort key for the row
func (c *Car) Keys(gsi int) (string, string, error) {
	// For this example, assuming GUID is the partition key and Email is the sort key.
	// Additional logic can be added to handle different GSIs if necessary.
	switch gsi {
	case 0: // Primary keys
		return c.OwnerEmail, fmt.Sprintf("%s-%s-%d", c.Make, c.Model, c.Year), nil
	default:
		// Handle other GSIs or return an error
		return "", "", nil
	}
}
