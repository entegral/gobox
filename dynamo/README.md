# Dynamo Package Usage Guide

This guide will help you understand how to use the `Row`, `MonoLink`, `DiLink`, and `TriLink` types in the `dynamo` package. These types are designed to be embedded into your own structs, providing them with access to DynamoDB.

## Prerequisites

Ensure that the following environment variables are configured:

- `AWS_REGION`: The AWS region where your DynamoDB is located.
- `AWS_PROFILE`: The AWS profile to use for authentication.
- `TABLENAME`: The name of the DynamoDB table to interact with.

## Row

The `Row` type is a basic type that represents a row in DynamoDB. It can be embedded into your own struct like so:

```go
type MyStruct struct {
    dynamo.Row
    // Your fields here
}
```

When a user embeds the [`Row`] struct from the [`dynamo`] package into their own struct, they gain access to several methods that [`Row`] provides. Here's a brief explanation of these methods:

- `Type() string`: This method returns the type of the record. By default, it returns "dilink" if the [`UnmarshalledType`] field is empty. Otherwise, it returns the value of [`UnmarshalledType`]. This method can be overridden by implementing it on the parent struct.

- `IsType(t string) bool`: This method checks if the record is of the given type `t`. It returns true if the type of the record matches `t`, and false otherwise.

- `TableName(ctx context.Context) string`: This method returns the name of the DynamoDB table. By default, it returns the value of the "TABLENAME" environment variable. If this environment variable is not set, it panics. This method can be overridden by implementing it on the parent struct.

- `Keys(gsi int) (partitionKey, sortKey string)`: This method is intended to return the partition key and sort key for the given Global Secondary Index (GSI). However, in the [`Row`] struct, this method simply panics with "not implemented". This method should be implemented in the parent struct to provide the correct keys.

Remember, these methods are available to the parent struct because of the way Go handles embedded structs. When you embed a struct in Go, the methods of the embedded struct become methods of the parent struct, and they can be overridden by implementing them on the parent struct.

You can find these methods in the [`Row`] struct in the [`dynamo`] package.

## MonoLink
The MonoLink type is used to establish a one-to-one relationship between two entities in DynamoDB. It can be embedded into your struct like so:

```go
package mypackage

import (
    "context"
    "github.com/entegral/gobox/dynamo"
    "github.com/entegral/gobox/types"
)

type MyEntity struct {
    types.Keyable
    // Your fields here
}

type MyStruct struct {
    dynamo.MonoLink[MyEntity]
    // Your fields here
}
```

In this example, MyEntity is a struct that implements the types.Linkable interface, which is required for the MonoLink struct.

Once you've embedded the MonoLink struct, you can use its methods. Here are some examples:

```go
package mypackage

import (
    "context"
    "fmt"
    "github.com/entegral/gobox/dynamo"
    "github.com/entegral/gobox/types"
)

func ExampleUsage() {
    // Create a new instance of MyStruct
    myStruct := MyStruct{
        MonoLink: dynamo.NewMonoLink(MyEntity{/* initialize your entity here */}),
        // Initialize your fields here
    }

    // Use the Link method to establish a connection between the entities
    err := myStruct.Link(context.Background(), myStruct)
    if err != nil {
        fmt.Println("Failed to link:", err)
        return
    }

    // Use the Unlink method to remove the connection between the entities
    err = myStruct.Unlink(context.Background(), myStruct)
    if err != nil {
        fmt.Println("Failed to unlink:", err)
        return
    }
}
```

In this example, the Link method is used to establish a connection between the entities, and the Unlink method is used to remove the connection. Note that these methods require a context and an instance of your entity. You would replace MyEntity{/* initialize your entity here */} with an actual instance of your entity.

Remember to handle errors appropriately in your own code. The error handling in this example is very basic and is just for demonstration purposes.

You can find these methods in the MonoLink struct in the dynamo package.

---
# DiLink

The DiLink type is used to establish a one-to-one relationship between two entities in DynamoDB. It can be embedded into your struct like so:

```go
package mypackage

import (
    "context"
    "github.com/entegral/gobox/dynamo"
    "github.com/entegral/gobox/types"
)

type MyEntity1 struct {
    types.Keyable
    // Your fields here
}

type MyEntity2 struct {
    types.Keyable
    // Your fields here
}

type MyStruct struct {
    dynamo.DiLink[MyEntity1, MyEntity2]
    // Your fields here
}
```

In this example, MyEntity1 and MyEntity2 are structs that implement the types.Linkable interface, which is required for the DiLink struct.

Once you've embedded the DiLink struct, you can use its methods. Here are some examples:

```go
package mypackage

import (
    "context"
    "fmt"
    "github.com/entegral/gobox/dynamo"
    "github.com/entegral/gobox/types"
)

func ExampleUsage() {
    // Create a new instance of MyStruct
    myStruct := MyStruct{
        DiLink: dynamo.NewDiLink(MyEntity1{/* initialize your entity here */}, MyEntity2{/* initialize your entity here */}),
        // Initialize your fields here
    }

    // Use the Link method to establish a connection between the entities
    err := myStruct.Link(context.Background(), MyEntity1{/* initialize your entity here */})
    if err != nil {
        fmt.Println("Failed to link:", err)
        return
    }

    // Use the Unlink method to remove the connection between the entities
    err = myStruct.Unlink(context.Background(), MyEntity1{/* initialize your entity here */})
    if err != nil {
        fmt.Println("Failed to unlink:", err)
        return
    }
}
```

## Finding Linked Objects

To use the [`DiLink`] struct to instantiate with one entity and load the other one, you would first need to create an instance of [`DiLink`] with one entity. Then, you can use the `LoadEntity1` method to load the other entity. Here's an example:

```go
package mypackage

import (
	"context"
	"fmt"
	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/types"
)

type MyEntity1 struct {
	types.Keyable
	// Your fields here
}

type MyEntity2 struct {
	types.Keyable
	// Your fields here
}

func ExampleUsage() {
	// Create a new instance of DiLink with one entity
	myDiLink := dynamo.CheckDiLink(MyEntity1{/* initialize your entity here */}, nil)

	// Use the LoadEntity1 method to load the other entity
	loaded, err := myDiLink.LoadEntity1(context.Background())
	if err != nil {
		fmt.Println("Failed to load entity:", err)
		return
	}

	if loaded {
		fmt.Println("Entity loaded successfully")
	} else {
		fmt.Println("Entity not found")
	}
}
```

In this example, `MyEntity1` and `MyEntity2` are structs that implement the `types.Linkable` interface, which is required for the [`DiLink`] struct. You would replace `MyEntity1{/* initialize your entity here */}` with an actual instance of your entity.

Remember to handle errors appropriately in your own code. The error handling in this example is very basic and is just for demonstration purposes.

You can find the `LoadEntity1` method in the [`DiLink`] struct in the [`dynamo`] package.


# TriLink
The TriLink type is used to establish a one-to-one relationship between three entities in DynamoDB. It can be embedded into your struct like so:

```go
type MyStruct struct {
    dynamo.TriLink[OtherStruct1, OtherStruct2, OtherStruct3]
    // Your fields here
}
```

You can then use the Link and Unlink methods to establish or remove the connection between the entities.