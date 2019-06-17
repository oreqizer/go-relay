package relay_test

import (
	"github.com/oreqizer/go-relay"
	"testing"
)

type Node struct {
	LocalID string
}

func (n *Node) ID() string {
	return n.LocalID
}

var (
	ONESTR   = "1"
	TWOSTR   = "2"
	THREESTR = "3"
	FOURSTR  = "4"
	FIVESTR  = "5"
)

var (
	ZERO = 0
	ONE  = 1
	TWO  = 2
	// THREE = 3
	FOUR = 4
	// FIVE  = 5
	SIX = 6
)

var nodes = []relay.Node{
	&Node{LocalID: ONESTR},
	&Node{LocalID: TWOSTR},
	&Node{LocalID: THREESTR},
	&Node{LocalID: FOURSTR},
	&Node{LocalID: FIVESTR},
}

var edges = []*relay.Edge{
	{Node: &Node{LocalID: ONESTR}, Cursor: ONESTR},
	{Node: &Node{LocalID: TWOSTR}, Cursor: TWOSTR},
	{Node: &Node{LocalID: THREESTR}, Cursor: THREESTR},
	{Node: &Node{LocalID: FOURSTR}, Cursor: FOURSTR},
	{Node: &Node{LocalID: FIVESTR}, Cursor: FIVESTR},
}

var tableConnectionFromArray = []struct {
	nodes []relay.Node
	args  *relay.ConnectionArgs
	out   *relay.Connection
}{
	{
		nodes: []relay.Node{},
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: &ZERO, Last: &ZERO},
		out: &relay.Connection{
			Edges:    []*relay.Edge{},
			PageInfo: relay.PageInfo{HasNextPage: false, HasPreviousPage: false, StartCursor: nil, EndCursor: nil},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: &ZERO, Last: &ZERO},
		out: &relay.Connection{
			Edges:    edges,
			PageInfo: relay.PageInfo{HasNextPage: false, HasPreviousPage: false, StartCursor: &ONESTR, EndCursor: &FIVESTR},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: &TWO, Last: &ZERO},
		out: &relay.Connection{
			Edges:    edges[:2],
			PageInfo: relay.PageInfo{HasNextPage: true, HasPreviousPage: false, StartCursor: &ONESTR, EndCursor: &TWOSTR},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: nil, After: nil, First: &ZERO, Last: &TWO},
		out: &relay.Connection{
			Edges:    edges[3:],
			PageInfo: relay.PageInfo{HasNextPage: false, HasPreviousPage: true, StartCursor: &FOURSTR, EndCursor: &FIVESTR},
		},
	},
	{
		nodes: nodes,
		args:  &relay.ConnectionArgs{Before: &FOURSTR, After: &TWOSTR, First: &ZERO, Last: &ZERO},
		out: &relay.Connection{
			Edges:    edges[1:4],
			PageInfo: relay.PageInfo{HasNextPage: false, HasPreviousPage: false, StartCursor: &TWOSTR, EndCursor: &FOURSTR},
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

		if out.PageInfo.StartCursor != nil || e.out.PageInfo.StartCursor != nil {
			if *out.PageInfo.StartCursor != *e.out.PageInfo.StartCursor {
				t.Errorf("%d: Start cursor: got %v, want %v", i, *out.PageInfo.StartCursor, *e.out.PageInfo.StartCursor)
			}
		}

		if out.PageInfo.EndCursor != nil || e.out.PageInfo.EndCursor != nil {
			if *out.PageInfo.EndCursor != *e.out.PageInfo.EndCursor {
				t.Errorf("%d: End cursor: got %v, want %v", i, *out.PageInfo.EndCursor, *e.out.PageInfo.EndCursor)
			}
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
	first  *int
	last   *int
	out    *relay.Connection
}{
	{edges: []*relay.Edge{}, before: nil, after: nil, first: &ZERO, last: &ZERO, out: tableConnectionFromArray[0].out},
	{edges: edges, before: nil, after: nil, first: &ZERO, last: &ZERO, out: tableConnectionFromArray[1].out},
	{edges: edges, before: nil, after: nil, first: &TWO, last: &ZERO, out: tableConnectionFromArray[2].out},
	{edges: edges, before: nil, after: nil, first: &ZERO, last: &TWO, out: tableConnectionFromArray[3].out},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, first: &ZERO, last: &ZERO, out: tableConnectionFromArray[4].out},
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

		if out.PageInfo.StartCursor != nil || e.out.PageInfo.StartCursor != nil {
			if *out.PageInfo.StartCursor != *e.out.PageInfo.StartCursor {
				t.Errorf("%d: Start cursor: got %v, want %v", i, *out.PageInfo.StartCursor, *e.out.PageInfo.StartCursor)
			}
		}

		if out.PageInfo.EndCursor != nil || e.out.PageInfo.EndCursor != nil {
			if *out.PageInfo.EndCursor != *e.out.PageInfo.EndCursor {
				t.Errorf("%d: End cursor: got %v, want %v", i, *out.PageInfo.EndCursor, *e.out.PageInfo.EndCursor)
			}
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
	{edges: []*relay.Edge{}, before: nil, after: nil, len: 0},
	{edges: edges, before: nil, after: nil, len: 5},
	{edges: edges, before: nil, after: &TWOSTR, len: 4},
	{edges: edges, before: &FOURSTR, after: nil, len: 4},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, len: 3},
}

func TestApplyCursorsToEdges(t *testing.T) {
	for i, e := range tableApplyCursorsToEdges {
		out := relay.ApplyCursorsToEdges(e.edges, e.before, e.after)
		if len(out) != e.len {
			t.Errorf("%d: Length: got %d, want %d", i, len(out), e.len)
		}

		if len(out) > 0 {
			if cursor := out[len(out)-1].Cursor; e.before != nil && cursor != *e.before {
				t.Errorf("%d: Before: got %s, want %s", i, cursor, *e.before)
			}

			if cursor := out[0].Cursor; e.after != nil && cursor != *e.after {
				t.Errorf("%d: After: got %s, want %s", i, cursor, *e.after)
			}
		}
	}
}

var tableHasPreviousPage = []struct {
	edges  []*relay.Edge
	before *string
	after  *string
	last   *int
	out    bool
}{
	{edges: []*relay.Edge{}, before: nil, after: nil, last: nil, out: false},
	{edges: edges, before: nil, after: nil, last: nil, out: false},
	{edges: edges, before: nil, after: nil, last: &ZERO, out: false},
	{edges: edges, before: nil, after: nil, last: &SIX, out: false},
	{edges: edges, before: nil, after: nil, last: &FOUR, out: true},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, last: &FOUR, out: false},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, last: &ONE, out: true},
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
	first  *int
	out    bool
}{
	{edges: []*relay.Edge{}, before: nil, after: nil, first: nil, out: false},
	{edges: edges, before: nil, after: nil, first: nil, out: false},
	{edges: edges, before: nil, after: nil, first: &ZERO, out: false},
	{edges: edges, before: nil, after: nil, first: &SIX, out: false},
	{edges: edges, before: nil, after: nil, first: &FOUR, out: true},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, first: &FOUR, out: false},
	{edges: edges, before: &FOURSTR, after: &TWOSTR, first: &ONE, out: true},
}

func TestHasNextPage(t *testing.T) {
	for i, e := range tableHasNextPage {
		out := relay.HasNextPage(e.edges, e.before, e.after, e.first)
		if out != e.out {
			t.Errorf("%d: got %v, want %v", i, out, e.out)
		}
	}
}
