package svgdoc

import "testing"

func TestGroupEndHandlesSelfClosedGroup(t *testing.T) {
	svg := `<g class="nodes"><g class="edgeLabels"/><g class="rough-node"><text>x</text></g></g>`
	end, ok := GroupEnd(svg, NodesMarker)
	if !ok {
		t.Fatal("GroupEnd did not find a closed node layer")
	}
	if end != len(svg) {
		t.Errorf("GroupEnd returned %d, want %d (the whole document)", end, len(svg))
	}
}

func TestGroupEndReportsMissingGroup(t *testing.T) {
	if _, ok := GroupEnd(`<g class="clusters"></g>`, NodesMarker); ok {
		t.Error("GroupEnd reported an absent node layer as closed")
	}
}

func TestClusterTitlesPaintLast(t *testing.T) {
	for _, testCase := range []struct {
		name string
		svg  string
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
