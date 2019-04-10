# Relay

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/oreqizer/go-relay)
[![Build Status](https://travis-ci.org/oreqizer/go-relay.svg?branch=master)](https://travis-ci.org/oreqizer/go-relay)
[![codecov](https://codecov.io/gh/oreqizer/go-relay/branch/master/graph/badge.svg)](https://codecov.io/gh/oreqizer/go-relay)

Implementation-agnostic utility functions, structures and interfaces for building [Relay](https://facebook.github.io/relay/docs/en/graphql-server-specification.html) compliant **GraphQL** servers.

It ain't much but it's honest work.

## API

### IDs

There are two functions - `ToGlobalID` and `FromGlobalID`. They behave the same like the JS reference implementation.

```go
var global = relay.ToGlobalID("User", "asdf") // Returns a base64 encoded string

var local = relay.FromGlobalID(global) // local.Type == "User", local.ID == "asdf"
```

### Connections

The only function you care about is `ConnectionFromArray`.

Make your types satisfy the `Node` interface and create the `ConnectionArgs` object, feed it into it and you'll get a `Connection`.

```go
// Our model
type User struct {
	LocalID string
	Name    string
}

// Satisfy the Node interface
func (u *User) ID() string {
	return u.LocalID
}

var nodes = []relay.Node{
	&User{LocalID: "1", Name: "Lol"},
	&User{LocalID: "2", Name: "Kek"},
	&User{LocalID: "3", Name: "Bur"},
}

var args = &relay.ConnectionArgs{Before: nil, After: nil, First: 2, Last: nil}

var conn = relay.ConnectionFromArray(nodes, args) // There you go!
```

## License

MIT
