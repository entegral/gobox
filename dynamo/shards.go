package dynamo

import (
	"fmt"
	"math/rand"
	"strings"
)

type Shard struct {
	stringFormatter string
	maxShard        int

	// PkShard is a field that is used to link to all the rows that are part of the same shard.
	PkShard string `dynamodbav:"pkshard,omitempty" json:"pkshard,omitempty"`
}

func (s Shard) MaxShard() int {
	if s.maxShard == 0 {
		return 100
	}
	return s.maxShard
}

func (s Shard) SetMaxShard(max int) {
	s.maxShard = max
}

// SetStringFormatter validates the string formatter and sets it to the Shard.
// The string formatter must contain a %d value to be used to format the shard number.
func (s Shard) SetStringFormatter(formatter string) error {
	if formatter == "" {
		return fmt.Errorf("string formatter cannot be empty")
	}
	if !strings.Contains(formatter, "%d") {
		return fmt.Errorf("string formatter must contain a %%d value")
	}
	s.stringFormatter = formatter
	return nil
}

// GetShard returns the shard formatted string with a random shard number.
// The shard number is a random number between 0 and the maxShard value.
// The string formatter is used to format the shard number into the string,
// and and it expects a %d value to be passed to it.
func (s Shard) GetShard() string {
	if s.stringFormatter == "" {
		return fmt.Sprintf("shard-%d", rand.Intn(s.MaxShard()))
	}
	return fmt.Sprintf(s.stringFormatter, rand.Intn(s.MaxShard()))
}
