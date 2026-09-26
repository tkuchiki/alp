package version

import "runtime/debug"

// Resolve returns a linker-provided version or, when unavailable, the version
// embedded by the Go toolchain for a module build.
func Resolve(linked string) string {
	info, ok := debug.ReadBuildInfo()
	return resolve(linked, info, ok)
}

func resolve(linked string, info *debug.BuildInfo, ok bool) string {
	if linked != "" {
		return linked
	}

	if !ok || info.Main.Version == "(devel)" {
		return ""
	}

	return info.Main.Version
}
