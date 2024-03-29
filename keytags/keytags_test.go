package keytags

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCompositeKeys(t *testing.T) {
	t.Run("primary keys", func(t *testing.T) {
		// Define a struct to hold the key parts
		type User struct {
			ID         string `pk:"1,prepend=USER:"`
			Email      string `pk:"2,prepend=EMAIL:"`
			UserType   string `sk:"1,prepend=TYPE:"`
			TimestampS int64  `sk:"2,prepend=TS:"`
			// Other fields...
		}

		ts := time.Now().Unix()
		// Create a new User object
		u := User{
			Email:      "test@gmail.com",
			ID:         "123",
			UserType:   "admin",
			TimestampS: ts,
		}

		// Generate the composite keys
		keyMap, err := generateCompositeKeys(&u, "/")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "USER:123/EMAIL:test@gmail.com", keyMap["pk"])
		assert.Equal(t, fmt.Sprintf("TYPE:admin/TS:%d", ts), keyMap["sk"])
	})

	t.Run("primary keys", func(t *testing.T) {
		// Define a struct to hold the key parts
		type User struct {
			ID         string `pk:"1,prepend=USER:" pk1:"1,prepend=USER:"`
			Email      string `pk:"2,prepend=EMAIL:" pk2:"1,prepend=EMAIL:"`
			UserType   string `sk:"1,prepend=TYPE:"`
			TimestampS int64  `sk:"2,prepend=TS:"`
			EyeColor   string `sk1:"1,prepend=EYECOLOR:"`
			HairColor  string `sk2:"2,prepend=HAIRCOLOR:"`
		}

		ts := time.Now().Unix()
		// Create a new User object
		u := User{
			Email:      "test@gmail.com",
			ID:         "123",
			UserType:   "admin",
			TimestampS: ts,
			EyeColor:   "blue",
			HairColor:  "black",
		}

		// Generate the composite keys
		keyMap, err := generateCompositeKeys(&u, "/")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "USER:123/EMAIL:test@gmail.com", keyMap["pk"])
		assert.Equal(t, fmt.Sprintf("TYPE:admin/TS:%d", ts), keyMap["sk"])
		assert.Equal(t, "USER:123", keyMap["pk1"])
		assert.Equal(t, "EMAIL:test@gmail.com", keyMap["pk2"])
		assert.Equal(t, "EYECOLOR:blue", keyMap["sk1"])
		assert.Equal(t, "HAIRCOLOR:black", keyMap["sk2"])
	})

}

func TestGenerateCompositeKeysBasic(t *testing.T) {
	type User struct {
		ID    string `pk:"1,prepend=USER#"`
		Email string `pk:"2,append=#EMAIL"`
	}

	user := User{ID: "123", Email: "test@example.com"}
	expected := map[string]string{"pk": "USER#123/test@example.com#EMAIL"}

	keys, err := generateCompositeKeys(&user, "/")
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	if keys["pk"] != expected["pk"] {
		t.Errorf("Expected pk %s, got %s", expected["pk"], keys["pk"])
	}
}

func TestGenerateCompositeKeysEmptyField(t *testing.T) {
	type User struct {
		ID    string `pk:"1,prepend=USER#"`
		Email string `pk:"2,append=#EMAIL"`
	}

	user := User{ID: "", Email: "test@example.com"} // ID is empty

	_, err := generateCompositeKeys(&user, "#")
	if err == nil {
		t.Error("Expected an error for empty ID field, but got nil")
	}
}

func TestGenerateCompositeKeysMixedTypes(t *testing.T) {
	type Order struct {
		OrderID  int    `pk:"1,prepend=ORDER#"`
		Customer string `pk:"2,append=#CUSTOMER"`
	}

	order := Order{OrderID: 456, Customer: "JohnDoe"}

	expected := map[string]string{"pk": "ORDER#456#JohnDoe#CUSTOMER"}

	keys, err := generateCompositeKeys(&order, "#")
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	if keys["pk"] != expected["pk"] {
		t.Errorf("Expected pk %s, got %s", expected["pk"], keys["pk"])
	}
}

func TestGenerateCompositeKeysMultipleKeyTypes(t *testing.T) {
	type Product struct {
		ProductID string `pk:"1,prepend=PROD:" sk:"1,prepend=SKU:"`
		SKU       string `sk:"2,append=:STOCK"`
	}

	product := Product{ProductID: "789", SKU: "XYZ123"}
	expected := map[string]string{
		"pk": "PROD:789",
		"sk": "SKU:789-XYZ123:STOCK",
	}

	keys, err := generateCompositeKeys(&product, "-")
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	if keys["pk"] != expected["pk"] || keys["sk"] != expected["sk"] {
		t.Errorf("Expected keys %v, got %v", expected, keys)
	}
}

func TestGenerateCompositeKeysZeroValuePrimaryKey(t *testing.T) {
	type Product struct {
		ID   int    `pk:"1,prepend=ID#"`
		Name string `pk:"2,append=#NAME"`
	}

	// ID has a zero value.
	product := Product{ID: 0, Name: "TestProduct"}

	_, err := generateCompositeKeys(&product, ":")
	if err == nil {
		t.Error("Expected an error for zero value in primary key field, but got nil")
	} else {
		t.Log("Received expected error:", err)
	}
}

func TestGenerateCompositeKeysNonPrimitiveField(t *testing.T) {
	type Address struct {
		City string
	}
	type User struct {
		Name    string  `pk:"1,prepend=NAME#"`
		Address Address `pk:"2,append=#ADDR"` // Address is a non-primitive type.
	}

	user := User{Name: "JohnDoe", Address: Address{City: "Springfield"}}

	_, err := generateCompositeKeys(&user, ":")
	if err == nil {
		t.Error("Expected an error for non-primitive type field used for key generation, but got nil")
	} else {
		t.Log("Received expected error:", err)
	}
}

func TestGenerateCompositeKeysSliceField(t *testing.T) {
	type Inventory struct {
		ProductID string   `pk:"1,prepend=PRODUCT#"`
		Items     []string `pk:"2,append=#ITEMS"` // Items is a slice.
	}

	inventory := Inventory{ProductID: "12345", Items: []string{"Item1", "Item2"}}

	_, err := generateCompositeKeys(&inventory, ":")
	if err == nil {
		t.Error("Expected an error for slice field used for key generation, but got nil")
	} else {
		t.Log("Received expected error:", err)
	}
}
