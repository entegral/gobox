# Dynamo Package Usage Guide

This guide provides insights on utilizing the `Row`, `MonoLink`, `DiLink`, and `TriLink` types in the dynamo package. These types, designed for embedding into your structs, grant access to DynamoDB functionalities.

## Prerequisites

Before using the dynamo package, ensure these environment variables are set:

- `AWS_REGION`: The AWS region of your DynamoDB instance.
- `AWS_PROFILE`: The AWS profile for authentication.
- `TABLENAME`: The name of the DynamoDB table for interactions.

## Row

The `Row` type represents a DynamoDB row and is the foundational type. Embed it into your struct as follows:

```gp
type MyStruct struct {
    dynamo.Row
    // Additional fields
}
```

### Key Methods of `Row`

- `Type() string`: Returns the record's type, defaulting to "dilink" or the value of `UnmarshalledType`.
- `IsType(t string) bool`: Checks if the record type matches the specified type `t`.
- `TableName(ctx context.Context) string`: Retrieves the DynamoDB table name, defaulting to the "TABLENAME" environment variable.
- `Keys(gsi int) (partitionKey, sortKey string)`: Intended to return keys for a Global Secondary Index (GSI). This method should be implemented in the parent struct.

These methods become part of the parent struct due to Go's embedded struct behavior, allowing for overriding in the parent struct.

## MonoLink

The `MonoLink` type establishes a one-to-one relationship between entities in DynamoDB. Embed it as follows:

```go
type MyStruct struct {
    dynamo.MonoLink[MyEntity]
    // Additional fields
}
```

### Using `MonoLink`

- `Link` method establishes a connection between entities.
- `Unlink` method removes the connection.

Example:

```go
func ExampleUsage() {
    myStruct := MyStruct{
        MonoLink: dynamo.NewMonoLink(MyEntity{/* initialize your entity here */}),
        // Additional initialization
    }

    // Linking and unlinking entities
    // Handle errors appropriately
}
```

## DiLink

`DiLink` creates a one-to-one relationship between two entities. Embed it as follows:

```go
type MyStruct struct {
    dynamo.DiLink[MyEntity1, MyEntity2]
    // Additional fields
}
```

### Using `DiLink`

- `Link` and `Unlink` methods manage connections between two entities.
- `LoadEntity1` loads the other entity in a `DiLink` relationship.

Example:

```go
func ExampleUsage() {
    myDiLink := dynamo.CheckDiLink(MyEntity1{/* initialize entity */}, nil)

    // Loading entities and error handling
}
```

## TriLink

`TriLink` is for establishing relationships among three entities in DynamoDB. Embed it in your struct:

```go
type MyStruct struct {
    dynamo.TriLink[OtherStruct1, OtherStruct2, OtherStruct3]
    // Additional fields
}
```

Utilize the `Link` and `Unlink` methods to establish or remove connections between these entities. 

Remember to handle errors and implement required methods in your structs to fully utilize the dynamo package's capabilities.
