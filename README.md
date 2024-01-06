# gobox
Robert Bruce

## Introduction

The purpose of this gobox is to provide a simple way to rapidly create new type directories for an application. Nearly all of my projects have a type directory, and many times those types need to be stored in a database. This gobox provides a simple framework for creating new types that have DynamoDB methods built into them, as well as a simple way to relate those types to each other in an easy and declarative way. 


## Overview

The bulk of the code is in the `dynamo` package. This package provides a simple way to create new types that have DynamoDB methods built into them. The `dynamo` package also provides a simple way to relate those types to each other in an easy and declarative way.

### Row Type

The `Row` type is the base type for all types that will be stored in DynamoDB. It provides a simple way to create new types that have DynamoDB methods built into them. It can be [embedded](examples/exampleLib/user.go) into other types [to provide basic CRUD operations](<examples/row/put/main.go>) on any data type.

These types are very important and should represent the core data types of your application. 

### Link Types

There are 3 `Link` types that this package provides that allow you to relate `Row` types to each other in various ways. These types are the `MonoLink`, `DiLink`, and `TriLink` types.

#### MonoLink 

The `MonoLink` type can be thought of as an [extension](examples/exampleLib/user.go) of the base `Row` type. It allows you to isolate less-frequently accessed fields of a `Row` type into a separate partition of the database by specifying an `Entity0` base entity. Since the `MonoLink`'s keys are deterministically derived from the `Row` type's keys, it is easy to access the `MonoLink`'s data from the base `Row` type. However, it is unlikely that the `MonoLink` would act as an entrypoint to the `Row` type's data. While certainly possible to query for a `MonoLink`'s data directly, then load the base type from it, this is not the intended use case for a `MonoLink` and it is recommended that you only use the `MonoLink` as a cost-saving measure for less-frequently accessed fields of a `Row` type.

#### DiLink

The [DiLink](examples/exampleLib/pinkslip.go) type is similar to a `MonoLink`, but it also has an `Entity1` base entity, allowing you to relate two `Row` types together. It enables a many-to-many relationship between the two types by linking two instances of the `Row` types together. 

The existence of a `DiLink` row in dynamo indicates that the two `Row` types are related to each other. Like the `MonoLink`, the `DiLink`'s keys are deterministically derived from the two entities' keys methods, so it is easy to access the `DiLink`'s data from either of the `Row` types. The `DiLink` even has helper functions allowing you to load all possible `Entity0s` or `Entity1s`. 


#### TriLink

The `TriLink` type is similar to a `DiLink`, but it also has an `Entity2` base entity, allowing you to query for, and relate, three `Row` types together with a single row. While not as common of a use case as the `DiLink`, the `TriLink` is still a useful tool for relating three `Row` types together in situations where all three types are exist and could benefit from direct querying.


## Instructions 

This gobox is intended to be used as a library. It is not intended to be used as a standalone application. There are several examples in the `examples` directory that show various ways to use this gobox tool in your application.