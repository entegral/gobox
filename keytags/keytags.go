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

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // This gets the field's metadata.
		fieldValue := val.Field(i)   // This gets the field's value as reflect.Value.

		// Use fieldValue.Kind() to get the kind of the field's value.
		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			return nil, fmt.Errorf("slice or array field %s.%s cannot be used for key generation", typeName, field.Name)
		}

		// Check if the field's type is not one of the supported primitive types for key generation.
		if !fieldValue.Type().Comparable() || fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Map || fieldValue.Kind() == reflect.Func || fieldValue.Kind() == reflect.Chan {
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

			operationAndValue := strings.SplitN(parts[1], "=", 2) // Split operation from value.
			if len(operationAndValue) != 2 {
				return nil, fmt.Errorf("malformed tag for %s.%s: missing operation or value", typeName, field.Name)
			}

			operation := operationAndValue[0]
			value := operationAndValue[1]

			fieldValue := fmt.Sprintf("%v", val.Field(i).Interface())
			if (fieldValue == "0" || strings.TrimSpace(fieldValue) == "" || fieldValue == "false") && (keyName == "pk" || strings.HasPrefix(keyName, "pk")) {
				return nil, fmt.Errorf("zero or empty value for field %s.%s used in %s", typeName, field.Name, keyName)
			}

			var keyValue string
			if operation == "prepend" {
				keyValue = value + fieldValue
			} else if operation == "append" {
				keyValue = fieldValue + value
			} else {
				keyValue = fieldValue // Just use the field value if no valid operation is specified.
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
