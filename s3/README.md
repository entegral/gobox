# BucketManager

The `BucketManager` is a struct that can be embedded into other structs to provide S3 functionality. It simplifies the process of interacting with S3 by providing methods to put, get, and delete objects, especially when using `Linkable` items from the `gobox` package.

## Usage

First, import the package:

```go
import "github.com/entegral/gobox/s3"
```

Next, define a struct that embeds the `BucketManager`:

```go
type MyStruct struct {
  dynamo.Row      // Row is Keyable, and provides an s3-friendly key
  *s3.BucketManager
  // other fields
  // ...
}
```

You can then create a new instance of your struct with a `BucketManager`:

```go
item := &MyStruct{
  BucketManager: s3.NewBucketManager("my-bucket"),

}
```

### PutObject

To put an object into the bucket, use the `PutObject` method:

```go
err := item.PutObject(context.Background(), item)
if err != nil {
  log.Fatalf("failed to put object: %v", err)
}
```

### GetObject

To get an object from the bucket, use the `GetObject` method:

```go
err := item.GetObject(context.Background(), item)
if err != nil {
  log.Fatalf("failed to get object: %v", err)
}
```

### DeleteObject

To delete an object from the bucket, use the `DeleteObject` method:

```go
err := item.DeleteObject(context.Background(), item)
if err != nil {
  log.Fatalf("failed to delete object: %v", err)
}
```
