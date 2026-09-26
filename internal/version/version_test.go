package version

import (
	"runtime/debug"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name   string
		linked string
		info   *debug.BuildInfo
		ok     bool
		want   string
	}{
		{name: "linked version", linked: "v1.2.3", want: "v1.2.3"},
		{name: "module version", info: &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}}, ok: true, want: "v1.2.3"},
		{name: "development build", info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, ok: true},
		{name: "build info unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.linked, tt.info, tt.ok); got != tt.want {
				t.Fatalf("resolve() = %q, want %q", got, tt.want)
			}
		})
	}
}
