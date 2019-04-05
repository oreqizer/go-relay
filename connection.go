package relay

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
	Edges    []*Edge   `json:"edges"`
	PageInfo *PageInfo `json:"pageInfo"`
}

/*
ConnectionArgs holds information about a connection arguments

https://facebook.github.io/relay/graphql/connections.htm#sec-Reserved-Types
*/
type ConnectionArgs struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
	First  int     `json:"first"`
	Last   int     `json:"last"`
}

/*
ConnectionFromSlice creates a connection from a slice of nodes
*/
func ConnectionFromSlice(nodes []Node, args *ConnectionArgs) *Connection {
	edges := make([]*Edge, len(nodes))
	for _, n := range nodes {
		edges = append(edges, &Edge{
			Node:   n,
			Cursor: n.ID(),
		})
	}

	return EdgesToReturn(edges, args.Before, args.After, args.First, args.Last)
}

/*
EdgesToReturn TODO

Consider returning an error like in
https://facebook.github.io/relay/graphql/connections.htm#sec-Pagination-algorithm
*/
func EdgesToReturn(all []*Edge, before, after *string, first, last int) *Connection {
	edges := ApplyCursorsToEdges(all, before, after)

	if first > 0 && first < len(edges) {
		edges = edges[:first]
	}

	if last > 0 && last < len(edges) {
		edges = edges[:last]
	}

	var startCursor, endCursor *string
	if fst := edges[0]; fst != nil {
		str := fst.Cursor
		startCursor = &str
	}

	if lst := edges[len(edges)-1]; lst != nil {
		str := lst.Cursor
		endCursor = &str
	}

	return &Connection{
		Edges: edges,
		PageInfo: &PageInfo{
			HasPreviousPage: HasPreviousPage(all, before, after, last),
			HasNextPage:     HasNextPage(all, before, after, first),
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

/*
ApplyCursorsToEdges TODO
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
				edges = edges[:i]
				break
			}
		}
	}
	return edges
}

/*
HasPreviousPage TODO
*/
func HasPreviousPage(all []*Edge, before, after *string, last int) bool {
	if last > 0 {
		edges := ApplyCursorsToEdges(all, before, after)
		return len(edges) > last
	}

	if after != nil {
		last := all[len(all)-1]
		return last.Cursor != *after
	}

	return false
}

/*
HasNextPage TODO
*/
func HasNextPage(all []*Edge, before, after *string, first int) bool {
	if first > 0 {
		edges := ApplyCursorsToEdges(all, before, after)
		return len(edges) > first
	}

	if before != nil {
		first := all[0]
		return first.Cursor != *before
	}

	return false
}
