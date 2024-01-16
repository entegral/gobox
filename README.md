# gobox
**Author:** Robert Bruce

## Introduction

gobox is designed to streamline the creation of new type directories for applications. Most projects necessitate a type directory, often requiring database storage. gobox offers an efficient framework for generating new types equipped with DynamoDB methods, facilitating straightforward and declarative inter-type relationships. 

## Overview

The core functionality resides in the `dynamo` package, which simplifies the creation of new types integrated with DynamoDB methods. It also enables easy, declarative relationships between these types.

### Row Type

At the foundation is the `Row` type, essential for all types destined for DynamoDB storage. It enables the creation of new types with inherent DynamoDB methods. By embedding the `Row` type into other types, it provides basic CRUD operations on any data type, representing the essential data types of your application.

### Link Types

To interrelate `Row` types, the package offers three `Link` types: `MonoLink`, `DiLink`, and `TriLink`.

#### MonoLink 

`MonoLink`, an augmentation of the `Row` type, segregates seldom-accessed `Row` fields into a separate database partition, using an `Entity0` base entity. With keys derived from the `Row` type’s keys, accessing `MonoLink` data is straightforward. However, `MonoLink` is not typically the primary access point for `Row` data. While direct querying of `MonoLink` data is feasible, it is primarily a cost-effective solution for infrequently accessed `Row` fields.

#### DiLink

`DiLink` expands on `MonoLink` by adding an `Entity1` base entity, linking two `Row` types. It fosters a many-to-many relationship, interconnecting two `Row` instances. The presence of a `DiLink` row signifies a relationship between two `Row` types. Its keys, derived from the entities’ keys, simplify access to `DiLink` data from either `Row` type, with helper functions for loading potential `Entity0s` or `Entity1s`. 

#### TriLink

`TriLink` is an extension of `DiLink`, incorporating an `Entity2` base entity. This allows for querying and relating three `Row` types with a single row. Though less common than `DiLink`, `TriLink` is valuable for interrelating three `Row` types in scenarios where such a relationship is beneficial.

## Instructions 

gobox is designed to be utilized as a library component within your projects, rather than as a standalone application. To guide you in implementing gobox's features, we have provided a comprehensive suite of tests. These tests not only assure the functionality of the library but also serve as practical examples. 

In the tests, you'll find various types implementing the tools and functionalities offered by gobox. These examples illustrate how to effectively use the `Row`, `MonoLink`, `DiLink`, and `TriLink` types, among other features. By studying these tests, you can gain insights into how to integrate and leverage the capabilities of gobox in your own applications.

We recommend examining these tests to familiarize yourself with the library's usage patterns and to see the implementation of its key components in action.

