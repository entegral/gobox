package keytags

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Define a struct to hold parts of the key along with their order.
type keyPart struct {
	Value string
	Order int
}

func generateCompositeKeys(v interface{}, delimiter string) (map[string]string, error) {
	val := reflect.ValueOf(v).Elem()
	typeName := val.Type().Name() // Get the type name of the struct.

	keys := make(map[string][]keyPart)
	orderUsed := make(map[string]map[int]bool) // New: Track orders used for each key type.

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i) // Ensure you're working with the field's value.

		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			return nil, fmt.Errorf("slice or array field %s.%s cannot be used for key generation", typeName, field.Name)
		}

		if !fieldValue.IsValid() || fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Map || fieldValue.Kind() == reflect.Func || fieldValue.Kind() == reflect.Chan {
			return nil, fmt.Errorf("non-primitive field %s.%s cannot be used for key generation", typeName, field.Name)
		}

		for _, keyName := range []string{"pk", "sk", "pk1", "sk1", "pk2", "sk2", "pk3", "sk3", "pk4", "sk4", "pk5", "sk5", "pk6", "sk6"} {
			tag, ok := field.Tag.Lookup(keyName)
			if !ok {
				continue // Skip if the tag is not present for the keyName.
			}

			parts := strings.Split(tag, ",")
			if len(parts) < 2 {
				continue // Ensure at least the order and one operation are specified.
			}

			order, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid order in %s tag for %s.%s: %v", keyName, typeName, field.Name, err)
			}

			// New: Check for duplicate orders within the same key type.
			if _, exists := orderUsed[keyName]; !exists {
				orderUsed[keyName] = make(map[int]bool)
			}
			if orderUsed[keyName][order] {
				return nil, fmt.Errorf("duplicate key ordering specified for %s in %s, both specify %d term", keyName, typeName, order)
			}
			orderUsed[keyName][order] = true

			operationAndValue := strings.SplitN(parts[1], "=", 2) // Split operation from value.
			if len(operationAndValue) != 2 {
				return nil, fmt.Errorf("malformed tag for %s.%s: missing operation or value", typeName, field.Name)
			}

			operation := operationAndValue[0]
			value := operationAndValue[1]

			// Corrected zero or empty value check
			if (keyName == "pk" || strings.HasPrefix(keyName, "pk")) && (fieldValue.Kind() == reflect.String && strings.TrimSpace(fieldValue.String()) == "" || !fieldValue.IsValid() || reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface())) {
				return nil, fmt.Errorf("zero or empty value for field %s.%s used in %s", typeName, field.Name, keyName)
			}

			var keyValue string
			if operation == "prepend" {
				keyValue = value + fmt.Sprint(fieldValue)
			} else if operation == "append" {
				keyValue = fmt.Sprint(fieldValue) + value
			} else {
				keyValue = fmt.Sprint(fieldValue) // Just use the field value if no valid operation is specified.
			}

			keys[keyName] = append(keys[keyName], keyPart{Value: keyValue, Order: order})
		}
	}

	finalKeys := make(map[string]string)
	for keyType, parts := range keys {
		sort.Slice(parts, func(i, j int) bool { return parts[i].Order < parts[j].Order })
		finalKeys[keyType] = concatenateKeyParts(parts, delimiter)
	}

	return finalKeys, nil
}

func concatenateKeyParts(parts []keyPart, delimiter string) string {
	var result []string
	for _, part := range parts {
		result = append(result, part.Value)
	}
	return strings.Join(result, delimiter)
}
