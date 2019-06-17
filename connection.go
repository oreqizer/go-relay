package relay

/*
Node is an interface for types satisfying Relay's Node interface

https://facebook.github.io/relay/graphql/objectidentification.htm#sec-Node-Interface
*/
type Node interface {
	ID() string
}

/*
Edge is an interface holding a node and a cursor

https://facebook.github.io/relay/graphql/connections.htm#sec-Edge-Types
*/
type Edge struct {
	Node   Node   `json:"node"`
	Cursor string `json:"cursor"`
}

/*
PageInfo holds information about connection edges

https://facebook.github.io/relay/graphql/connections.htm#sec-Reserved-Types
*/
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

/*
Connection holds information about a connection

https://facebook.github.io/relay/graphql/connections.htm#sec-Reserved-Types
*/
type Connection struct {
	Edges    []*Edge  `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}

/*
ConnectionArgs holds information about a connection arguments

https://facebook.github.io/relay/graphql/connections.htm#sec-Reserved-Types
*/
type ConnectionArgs struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
	First  *int    `json:"first"`
	Last   *int    `json:"last"`
}

/*
ConnectionFromArray creates a connection from an array of nodes
*/
func ConnectionFromArray(nodes []Node, args *ConnectionArgs) *Connection {
	if args == nil {
		return nil
	}

	edges := make([]*Edge, len(nodes))
	for i, n := range nodes {
		edges[i] = &Edge{
			Node:   n,
			Cursor: n.ID(),
		}
	}

	return EdgesToReturn(edges, args.Before, args.After, args.First, args.Last)
}

/*
EdgesToReturn slices edges according to arguments, returning a connection

Consider returning an error like in
https://facebook.github.io/relay/graphql/connections.htm#sec-Pagination-algorithm
*/
func EdgesToReturn(all []*Edge, before, after *string, first, last *int) *Connection {
	edges := ApplyCursorsToEdges(all, before, after)

	if first != nil && *first > 0 && *first < len(edges) {
		edges = edges[:*first]
	}

	if last != nil && *last > 0 && *last < len(edges) {
		edges = edges[len(edges)-*last:]
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		if fst := edges[0]; fst != nil {
			str := fst.Cursor
			startCursor = &str
		}

		if lst := edges[len(edges)-1]; lst != nil {
			str := lst.Cursor
			endCursor = &str
		}
	}

	return &Connection{
		Edges: edges,
		PageInfo: PageInfo{
			HasPreviousPage: HasPreviousPage(all, before, after, last),
			HasNextPage:     HasNextPage(all, before, after, first),
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

/*
ApplyCursorsToEdges slices edges according to cursors
*/
func ApplyCursorsToEdges(all []*Edge, before, after *string) []*Edge {
	edges := all
	if after != nil {
		for i, e := range edges {
			if e.Cursor == *after {
				edges = edges[i:]
				break
			}
		}
	}

	if before != nil {
		for i, e := range edges {
			if e.Cursor == *before {
				edges = edges[:i+1]
				break
			}
		}
	}
	return edges
}

/*
HasPreviousPage determines whether there's a previous page according to cursors

https://facebook.github.io/relay/graphql/connections.htm#sec-undefined.PageInfo.Fields
*/
func HasPreviousPage(all []*Edge, before, after *string, last *int) bool {
	if last != nil && *last > 0 {
		edges := ApplyCursorsToEdges(all, before, after)
		return len(edges) > *last
	}

	return false
}

/*
HasNextPage determines whether there's another page according to cursors

https://facebook.github.io/relay/graphql/connections.htm#sec-undefined.PageInfo.Fields
*/
func HasNextPage(all []*Edge, before, after *string, first *int) bool {
	if first != nil && *first > 0 {
		edges := ApplyCursorsToEdges(all, before, after)
		return len(edges) > *first
	}

	return false
}
