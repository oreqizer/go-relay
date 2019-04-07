# Relay

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/oreqizer/go-relay)
[![Build Status](https://travis-ci.org/oreqizer/go-relay.svg?branch=master)](https://travis-ci.org/oreqizer/go-relay)
[![codecov](https://codecov.io/gh/oreqizer/go-relay/branch/master/graph/badge.svg)](https://codecov.io/gh/oreqizer/go-relay)

Implementation-agnostic utility functions, structures and interfaces for building [Relay](https://facebook.github.io/relay/docs/en/graphql-server-specification.html) compliant **GraphQL** servers.

It ain't much but it's honest work.

## API

The only function you care about is `ConnectionFromArray`.

Make your types satisfy the `Node` interface and create the `ConnectionArgs` object, feed it into it and you'll get a `Connection`.

```go
// Our model
type User struct {
	Id string
	Name string
}

// Satisfy the Node interface
func (u *User) ID() string {
	return u.Id
}

var nodes = []relay.Node{
	&User{Id: "1", Name: "Lol"},
	&User{Id: "2", Name: "Kek"},
	&User{Id: "3", Name: "Bur"},
}

var args = &relay.ConnectionArgs{Before: nil, After: nil, First: 2, Last: 0}

var conn = relay.ConnectionFromArray(nodes, args) // There you go!
```

## License

MIT
