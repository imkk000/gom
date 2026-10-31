package cmd

import (
	"slices"
	"testing"

	"github.com/imkk000/gom/internal/bin"
)

func TestSelectPackages(t *testing.T) {
	paths := []string{"github.com/a/ccx", "github.com/b/lint/v2/cmd/golint", "github.com/c/tool/v3"}
	installed := map[string]bin.Binary{
		"github.com/c/tool/v3": {File: "/gobin/tool"},
	}
	cases := []struct {
		names []string
		want  []string
	}{
		{nil, paths},
		{[]string{"ccx"}, []string{"github.com/a/ccx"}},
		{[]string{"tool", "golint", "ccx", "tool"}, []string{"github.com/c/tool/v3", "github.com/b/lint/v2/cmd/golint", "github.com/a/ccx"}},
		{[]string{"github.com/c/tool/v3"}, []string{"github.com/c/tool/v3"}},
	}
	for _, c := range cases {
		got, err := selectPackages(paths, installed, c.names)
		if err != nil {
			t.Fatalf("selectPackages(%v): %v", c.names, err)
		}
		if !slices.Equal(got, c.want) {
			t.Errorf("selectPackages(%v) = %v; want %v", c.names, got, c.want)
		}
	}
	if _, err := selectPackages(paths, installed, []string{"missing"}); err == nil {
		t.Error("selectPackages(missing) returned no error")
	}
}
