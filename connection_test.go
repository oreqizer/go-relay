package relay_test

import (
	"github.com/oreqizer/go-relay"
	"testing"
)

type Node struct {
	Id string
}

func (n *Node) ID() string {
	return n.Id
}

var (
	ONE   = "1"
	TWO   = "2"
	THREE = "3"
	FOUR  = "4"
	FIVE  = "5"
)

var nodes = []relay.Node{
	&Node{Id: ONE},
	&Node{Id: TWO},
	&Node{Id: THREE},
	&Node{Id: FOUR},
	&Node{Id: FIVE},
}

var edges = []*relay.Edge{
	{Node: &Node{Id: ONE}, Cursor: ONE},
	{Node: &Node{Id: TWO}, Cursor: TWO},
	{Node: &Node{Id: THREE}, Cursor: THREE},
	{Node: &Node{Id: FOUR}, Cursor: FOUR},
	{Node: &Node{Id: FIVE}, Cursor: FIVE},
}

var tableConnectionFromArray = []struct {
	nodes []relay.Node
	args  *relay.ConnectionArgs
	out   *relay.Connection
}{
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: 0, Last: 0},
		out: &relay.Connection{
			Edges:    edges,
			PageInfo: &relay.PageInfo{HasNextPage: false, HasPreviousPage: false, StartCursor: &ONE, EndCursor: &FIVE},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: 2, Last: 0},
		out: &relay.Connection{
			Edges:    edges[:2],
			PageInfo: &relay.PageInfo{HasNextPage: true, HasPreviousPage: false, StartCursor: &ONE, EndCursor: &TWO},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: 0, Last: 2},
		out: &relay.Connection{
			Edges:    edges[3:],
			PageInfo: &relay.PageInfo{HasNextPage: false, HasPreviousPage: true, StartCursor: &FOUR, EndCursor: &FIVE},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: &FOUR, After: &TWO, First: 0, Last: 0},
		out: &relay.Connection{
			Edges:    edges[1:4],
			PageInfo: &relay.PageInfo{HasNextPage: false, HasPreviousPage: false, StartCursor: &TWO, EndCursor: &FOUR},
		},
	},
}

func TestConnectionFromArray(t *testing.T) {
	empty := relay.ConnectionFromArray(tableConnectionFromArray[0].nodes, nil)
	if empty != nil {
		t.Errorf("Expected nil output for nil args")
	}

	for i, e := range tableConnectionFromArray {
		out := relay.ConnectionFromArray(e.nodes, e.args)
		if out == nil {
			t.Errorf("Unexpected nil output")
			return
		}

		// PageInfo
		if out.PageInfo.HasNextPage != e.out.PageInfo.HasNextPage {
			t.Errorf("%d: Has next page: got %v, want %v", i, out.PageInfo.HasNextPage, e.out.PageInfo.HasNextPage)
		}

		if out.PageInfo.HasPreviousPage != e.out.PageInfo.HasPreviousPage {
			t.Errorf("%d: Has previous page: got %v, want %v", i, out.PageInfo.HasPreviousPage, e.out.PageInfo.HasPreviousPage)
		}

		if *out.PageInfo.StartCursor != *e.out.PageInfo.StartCursor {
			t.Errorf("%d: Start cursor: got %v, want %v", i, *out.PageInfo.StartCursor, *e.out.PageInfo.StartCursor)
		}

		if *out.PageInfo.EndCursor != *e.out.PageInfo.EndCursor {
			t.Errorf("%d: End cursor: got %v, want %v", i, *out.PageInfo.EndCursor, *e.out.PageInfo.EndCursor)
		}

		// Edges
		if len(out.Edges) != len(e.out.Edges) {
			t.Errorf("%d: Edges length: got %d, want %d", i, len(out.Edges), len(e.out.Edges))
			return
		}

		for j, eedge := range e.out.Edges {
			oedge := out.Edges[j]
			if eedge.Cursor != oedge.Cursor {
				t.Errorf("%d: Edge #%d: Cursor: got %s, want %s", i, j, oedge.Cursor, eedge.Cursor)
			}
		}
	}
}

var tableEdgesToReturn = []struct {
	edges  []*relay.Edge
	before *string
	after  *string
	first  int
	last   int
	out    *relay.Connection
}{
	{edges: edges, before: nil, after: nil, first: 0, last: 0, out: tableConnectionFromArray[0].out},
	{edges: edges, before: nil, after: nil, first: 2, last: 0, out: tableConnectionFromArray[1].out},
	{edges: edges, before: nil, after: nil, first: 0, last: 2, out: tableConnectionFromArray[2].out},
	{edges: edges, before: &FOUR, after: &TWO, first: 0, last: 0, out: tableConnectionFromArray[3].out},
}

func TestEdgesToReturn(t *testing.T) {
	for i, e := range tableEdgesToReturn {
		out := relay.EdgesToReturn(e.edges, e.before, e.after, e.first, e.last)
		// PageInfo
		if out.PageInfo.HasNextPage != e.out.PageInfo.HasNextPage {
			t.Errorf("%d: Has next page: got %v, want %v", i, out.PageInfo.HasNextPage, e.out.PageInfo.HasNextPage)
		}

		if out.PageInfo.HasPreviousPage != e.out.PageInfo.HasPreviousPage {
			t.Errorf("%d: Has previous page: got %v, want %v", i, out.PageInfo.HasPreviousPage, e.out.PageInfo.HasPreviousPage)
		}

		if *out.PageInfo.StartCursor != *e.out.PageInfo.StartCursor {
			t.Errorf("%d: Start cursor: got %v, want %v", i, *out.PageInfo.StartCursor, *e.out.PageInfo.StartCursor)
		}

		if *out.PageInfo.EndCursor != *e.out.PageInfo.EndCursor {
			t.Errorf("%d: End cursor: got %v, want %v", i, *out.PageInfo.EndCursor, *e.out.PageInfo.EndCursor)
		}

		// Edges
		if len(out.Edges) != len(e.out.Edges) {
			t.Errorf("%d: Edges length: got %d, want %d", i, len(out.Edges), len(e.out.Edges))
			return
		}

		for j, eedge := range e.out.Edges {
			oedge := out.Edges[j]
			if eedge.Cursor != oedge.Cursor {
				t.Errorf("%d: Edge #%d: Cursor: got %s, want %s", i, j, oedge.Cursor, eedge.Cursor)
			}
		}
	}
}

var tableApplyCursorsToEdges = []struct {
	edges  []*relay.Edge
	before *string
	after  *string
	len    int
}{
	{edges: edges, before: nil, after: nil, len: 5},
	{edges: edges, before: nil, after: &TWO, len: 4},
	{edges: edges, before: &FOUR, after: nil, len: 4},
	{edges: edges, before: &FOUR, after: &TWO, len: 3},
}

func TestApplyCursorsToEdges(t *testing.T) {
	for i, e := range tableApplyCursorsToEdges {
		out := relay.ApplyCursorsToEdges(e.edges, e.before, e.after)
		if len(out) != e.len {
			t.Errorf("%d: Length: got %d, want %d", i, len(out), e.len)
		}

		if cursor := out[len(out)-1].Cursor; e.before != nil && cursor != *e.before {
			t.Errorf("%d: Before: got %s, want %s", i, cursor, *e.before)
		}

		if cursor := out[0].Cursor; e.after != nil && cursor != *e.after {
			t.Errorf("%d: After: got %s, want %s", i, cursor, *e.after)
		}
	}
}

var tableHasPreviousPage = []struct {
	edges  []*relay.Edge
	before *string
	after  *string
	last   int
	out    bool
}{
	{edges: edges, before: nil, after: nil, last: 0, out: false},
	{edges: edges, before: nil, after: nil, last: 6, out: false},
	{edges: edges, before: nil, after: nil, last: 4, out: true},
	{edges: edges, before: &FOUR, after: &TWO, last: 4, out: false},
	{edges: edges, before: &FOUR, after: &TWO, last: 1, out: true},
}

func TestHasPreviousPage(t *testing.T) {
	for i, e := range tableHasPreviousPage {
		out := relay.HasPreviousPage(e.edges, e.before, e.after, e.last)
		if out != e.out {
			t.Errorf("%d: got %v, want %v", i, out, e.out)
		}
	}
}

var tableHasNextPage = []struct {
	edges  []*relay.Edge
	before *string
	after  *string
	first  int
	out    bool
}{
	{edges: edges, before: nil, after: nil, first: 0, out: false},
	{edges: edges, before: nil, after: nil, first: 6, out: false},
	{edges: edges, before: nil, after: nil, first: 4, out: true},
	{edges: edges, before: &FOUR, after: &TWO, first: 4, out: false},
	{edges: edges, before: &FOUR, after: &TWO, first: 1, out: true},
}

func TestHasNextPage(t *testing.T) {
	for i, e := range tableHasNextPage {
		out := relay.HasNextPage(e.edges, e.before, e.after, e.first)
		if out != e.out {
			t.Errorf("%d: got %v, want %v", i, out, e.out)
		}
	}
}
