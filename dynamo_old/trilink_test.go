package dynamo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Event will use the TriLink to link together a User, Venue, and Date
type Event struct {
	TriLink[*Person, *Venue, *Date]
	Name string `json:"name"`
}

func (e *Event) Type() string {
	return "event"
}

type Person struct {
	Row
	Name string `json:"name"`
}

func (p *Person) Type() string {
	return "person"
}

func (p *Person) Keys(gsi int) (string, string, error) {
	// for this example we will just use static keys
	return "personPk", "personSk", nil
}

type Venue struct {
	Row
	Name string `json:"name"`
}

func (v *Venue) Type() string {
	return "venue"
}

func (v *Venue) Keys(gsi int) (string, string, error) {
	// for this example we will just use static keys
	return "venuePk", "venueSk", nil
}

type Date struct {
	Row
	Date time.Time `json:"date"`
}

func (d *Date) Type() string {
	return "date"
}

func (d *Date) Keys(gsi int) (string, string, error) {
	// for this example we will just use static keys
	return "datePk", "dateSk", nil
}

func TestTriLink(t *testing.T) {
	ctx := context.Background()
	person := &Person{
		Name: "Example Person",
	}
	venue := &Venue{
		Name: "Example Venue",
	}
	date := &Date{
		Date: time.Now(),
	}
	exampleEvent := &Event{
		TriLink: *NewTriLink(person, venue, date), // it would be optimal to have a NewEvent(user, venue, date) method, but this is just an example
		Name:    "Example Event",
	}
	defer exampleEvent.Delete(ctx, exampleEvent)
	t.Run("CheckTriLink", func(t *testing.T) {
		person.Delete(ctx, person)
		venue.Delete(ctx, venue)
		date.Delete(ctx, date)
		exampleEvent.Delete(ctx, exampleEvent)
		t.Run("should return entity0 missing", func(t *testing.T) {
			linkExists, err := exampleEvent.CheckLink(ctx, exampleEvent, person, venue, date)
			if err == nil {
				t.Error("expected error")
			}
			assert.IsType(t, &ErrItemNotFound{}, err)
			assert.Equal(t, false, linkExists)
		})
		t.Run("should return entity1 missing", func(t *testing.T) {
			err := person.Put(ctx, person)
			if err != nil {
				t.Error(err)
			}
			linkExists, err := exampleEvent.CheckLink(ctx, exampleEvent, person, venue, date)
			if err == nil {
				t.Error("expected error")
			}
			assert.IsType(t, &ErrItemNotFound{}, err)
			assert.Equal(t, false, linkExists)
		})
		t.Run("should return entity2 missing", func(t *testing.T) {
			err := venue.Put(ctx, venue)
			if err != nil {
				t.Error(err)
			}
			linkExists, err := exampleEvent.CheckLink(ctx, exampleEvent, person, venue, date)
			if err == nil {
				t.Error("expected error")
			}
			assert.IsType(t, &ErrItemNotFound{}, err)
			assert.Equal(t, false, linkExists)
		})
		t.Run("should return false when the record isn't in dynamo", func(t *testing.T) {
			err := date.Put(ctx, date)
			if err != nil {
				t.Error(err)
			}
			linkExists, err := exampleEvent.CheckLink(ctx, exampleEvent, person, venue, date)
			assert.Nil(t, err)
			assert.Equal(t, false, linkExists)
		})
		t.Run("should return true when the record is in dynamo", func(t *testing.T) {
			err := exampleEvent.Put(ctx, exampleEvent)
			if err != nil {
				t.Error(err)
			}
			linkExists, err := exampleEvent.CheckLink(ctx, exampleEvent, person, venue, date)
			assert.Nil(t, err)
			assert.Equal(t, true, linkExists)
		})
	})

	t.Run("the TriLink operates like the DiLink, but it allows binding of up to 3 entities", func(t *testing.T) {

	})
}
