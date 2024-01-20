# Dynamo Package Usage Guide

This guide delves into the usage of the `Row`, `MonoLink`, `DiLink`, and `TriLink` types in the Dynamo package, designed for embedding into your structs to facilitate DynamoDB functionalities.


## Prerequisites

Before using the Dynamo package, ensure the following environment variables are set:

- `AWS_REGION`: Your DynamoDB instance's AWS region.
- `AWS_PROFILE`: AWS profile for authentication.
- `TABLENAME`: Name of your DynamoDB table.

Note: The client used by gobox will automatically load credentials from the IAM role assigned to the tool at runtime and it will use the value of the `TABLENAME` environment variable. Each type in the Dynamo package has the capacity to override this default table. If running locally, your `AWS_REGION` and an `AWS_PROFILE` are required.

## Row

`Row` represents a DynamoDB row and is foundational. It is a `Base` type and should contain all fields involved in your type's access patterns. Only fields that are involved in looking this item up should exist directly on this type. 

Satisfy the Typeable interface and you can already Get, Put and Delete these intstances to dynamo in a UUID-based, scalable way. Embed it as follows:

```go
type User struct {
  dynamo.Row
  // fields required to find this Base type in dynamo given your application's access patterns
  EMAIL string `json:"email"`
}

func (u *User) Type() string {
  return "User"
}


user := &User{
    EMAIL: "test@gmail.com"
}

func Example() {
    err := user.Put(context.TODO(), user) // this is the catch, you need to pass the instance into the method.
    if err != nil {
        // handle error
    }

    getUser := &User{}  // look ma, no fields!
    getUser.Pk = user.Pk // the Put method used the default Keyable method, which genrates a UUID if the Pk isn't set, so we need to set the Pk to the same value
    loaded, err := user.Get(context.TODO(), user)
    if err != nil {
        // handle error
    }
    // this is markdown, so we can't show the output, but getUser will be the same as the saved user
    fmt.Println(getUser)

    // guess what we do next...
    err = user.Delete(context.TODO(), user) // if you said "pass the instance to the method again", you're right!
    if err != nil {
        // handle error
    }

```

### Key Methods of Row

- `Keys(gsi int) (partitionKey, sortKey string, err error)`: Returns keys for a Composite Key or GSI.
- `Type() string`: Returns the record type, typically "dilink" or `UnmarshalledType`.
- `MaxShard() int`: Returns the shard count, defaulting to 100.
- `TableName(ctx context.Context) string`: Retrieves the DynamoDB table name.

## MonoLink

`MonoLink` establishes one-to-one relationships. Their keys are derrived from the keys of the `Base` type, so it is likely end up on a different dynamo partition when saved, allowing horizontally scalable groups of fields that relate to the base type. Embed it as follows:

```go
// UserDetails can be loaded independently of User by using User's keys as its own.
// This gets really cool when you define required

type UserDetails struct {
dynamo.MonoLink[*User]
  // examples of fields that shouldnt need to be every time:
  Name string `json:"name"`
  DateOfBirth time.Time `json:"dateOfBirth"`
  Address     string    `json:"address"`
  PhoneNumber string    `json:"phoneNumber"`
}

func (u *UserDetails) Type() string {
  return "UserDetails"
}

type UserPreferences struct {
dynamo.MonoLink[*User]
  // examples of fields that shouldnt need to be loaded every time:
  FavoriteColor string `json:"favoriteColor"`
  FavoriteFood  string `json:"favoriteFood"`
}

func (u *UserPreferences) Type() string {
  return "UserPreferences"
}
```

### Using MonoLink

All the xxxLink types embed the Row type too, so they can be used in the same way. The only difference is that they have a `Link` and `Unlink` method that can be used to manage the relationship between the Link and the Base type.

- `Get`: Retrieves the MonoLink, if it exists.
- `Link` (or `Put` if you like consistency): Saves the MonoLink.
- `Unlink` (or `Delete` if you like consistency): Removes the connection.

```go
func ExampleUsage() {
    user := &User{
        Pk: "HUMAN-READABLE-UNIQUE-UUID",
        EMAIL: "test@gmail.com",
    }


```

## DiLink
`DiLink` creates one-to-one relationships between two `Base` entities. It can also have its own fields, I'd limit those fields to things pertaining to the relationship yourself, but hey, you do you.

First we will define a 2nd `Base` type, then we will embed the `DiLink` type into it. Embed as follows:

```go

type Car struct {
    dynamo.Row
    Make string `json:"make"`
    Model string `json:"model"`
    Year int `json:"year"`
}

func (c *Car) Type() string {
    return "Car"
}
```

Now we will embed the `DiLink` type into a PinkSlip to establish ownership of a car. Embed as follows:

```go

type PinkSlip struct {
    dynamo.DiLink[*User*, *Car]
    DateOfSale time.Time `json:"dateOfSale"`
    Price int `json:"price"`
}
```

Of course you can also use the aformentioned methods for CRUD operations on instances of this type, but something extra cool happens when you put a DiLink, it actually writes the composite keys of the underlying types into specific fields:

```go
pinkSlip := &PinkSlip{
    DiLink: *dynamo.NewDiLink(user, car), // this is a little helper method that creates the DiLink for you. make sure the order of the types matches the order of the types in the struct.
    DateOfSale: time.Now(),
    Price: 10000,
}

err := pinkSlip.Put(context.TODO(), pinkSlip)
if err != nil {
    // handle error
}

if pinkSlip.E0pk != user.Pk {
    panic("well... ill eat my hat")
}

if pinkSlip.E1pk != car.Pk {
    panic("...guess ill eat my belt too!")
}

```

Why does this matter? Because there are helpers for loading links based on one of these entities. For example, if you want to load all the PinkSlips for a user, you can do this:

```go

pinkSlips, err := dynamo.FindLinksByEntity0[User, PinkSlip](ctx, user, "User")
if err != nil {
    // handle error
}

```

This will return all the PinkSlips for the user. You can also do the same thing for the other entity:

```go
pinkSlips, err := dynamo.FindLinksByEntity1[Car, PinkSlip](ctx, car, "Car")
if err != nil {
    // handle error
}

```

### Using DiLink

- `CheckLink`: Accepts two entity instances and attempts to laod them from dynamo. It will return an error if either of the entities do not exist. If the link does not exist, it will not return an error. Links should not be created if the entities do not exist.

    ```go
    func CheckLinkExample() {
        // lets assume we have already saved user and car, but not pinkSlip
        linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, &user, &car)
        if err != nil {
            return nil, err
        }
        if !linkExists {
            return pinkSlip, pinkSlip.Link(ctx, pinkSlip)
        } else {
            return pinkSlip, pinkSlip.Unlink(ctx, pinkSlip)
        }
    }

    ```

- `LoadEntity0`: This method will first attempt to call Keys(0) on entity0 and issue a dynamo.Get. If the entity has not been set with an entity0, it will attempt to extract the relevant keys from the DiLink's Composite keys and issue a dynamo.Get using those values. If an error occurs, it returns it, but if the entity is found, it returns nil and unmarhsals the entity into Entity0 field.

    ```go
    func LoadEntity0Example() {
        user := &User{
            Pk: "HUMAN-READABLE-UNIQUE-UUID",
        }
        pinkSlip.Entity0 = user

        // lets assume we have already saved user and car, but not pinkSlip
        err := pinkSlip.LoadEntity0(ctx, pinkSlip)
        if err != nil {
            return nil, err
        }
        // now the pinkSlip has this user loaded into the Entity0 field (assuming its in dynamo)
    }

    ```

- `LoadEntity1`: Basically the same as Entity0, but for Entity1.

## TriLink

The `TriLink` is used much in the same way as the `DiLink`, but links together three entities. Have fun!
