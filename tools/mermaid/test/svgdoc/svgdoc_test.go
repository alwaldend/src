package svgdoc

import "testing"

// TestGroupEndHandlesSelfClosedGroup checks that an empty child layer leaves the parent depth intact.
func TestGroupEndHandlesSelfClosedGroup(t *testing.T) {
	svg := `<g class="nodes"><g class="edgeLabels"/><g class="rough-node"><text>x</text></g></g>`
	end, ok := groupEnd(svg, nodesMarker)
	if !ok {
		t.Fatal("groupEnd did not find a closed node layer")
	}
	if end != len(svg) {
		t.Errorf("groupEnd returned %d, want %d (the whole document)", end, len(svg))
	}
}

// TestGroupEndReportsMissingGroup checks the missing-marker result.
func TestGroupEndReportsMissingGroup(t *testing.T) {
	if _, ok := groupEnd(`<g class="clusters"></g>`, nodesMarker); ok {
		t.Error("groupEnd reported an absent node layer as closed")
	}
}

// TestClusterTitlesPaintLast checks the ordering and presence of both SVG layers.
func TestClusterTitlesPaintLast(t *testing.T) {
	for _, testCase := range []struct {
		// name identifies the layering scenario.
		name string
		// svg contains the rendered layers in paint order.
		svg string
		// want reports whether the node layer precedes every title.
		want bool
	}{
		{
			name: "title lifted past the nodes",
			svg:  `<g class="nodes"></g><g class="cluster-label">t</g>`,
			want: true,
		},
		{
			name: "title still inside the clusters layer",
			svg:  `<g class="clusters"><g class="cluster-label">t</g></g><g class="nodes"></g>`,
			want: false,
		},
		{
			name: "no container title to order",
			svg:  `<g class="nodes"></g>`,
			want: false,
		},
		{
			name: "no node layer to compare against",
			svg:  `<g class="cluster-label">t</g>`,
			want: false,
		},
	} {
		if got := ClusterTitlesPaintLast(testCase.svg); got != testCase.want {
			t.Errorf("%s: ClusterTitlesPaintLast = %t, want %t",
				testCase.name, got, testCase.want)
		}
	}
}
