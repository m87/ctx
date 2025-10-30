package core

import "slices"

var (
	Release = "2.2.0"
	Commit  = "none"
	Date    = "unknown"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func Sort(versions []Version) {
	slices.SortFunc(versions, func(a, b Version) int {
		if a.Major != b.Major {
			return a.Major - b.Major
		}

		if a.Minor != b.Minor {
			return a.Minor - b.Minor
		}

		return a.Patch - b.Patch
	})
}
