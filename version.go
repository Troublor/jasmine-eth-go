package jasmine_eth_go

import "fmt"

const (
	major = 0
	minor = 0
	patch = 6
)

func VersionStr() string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}

func VersionNum() int {
	return major*1e6 + minor*1e3 + patch
}

func Version() struct {
	VersionStr string
	VersionNum int
} {
	return struct {
		VersionStr string
		VersionNum int
	}{
		VersionStr: VersionStr(),
		VersionNum: VersionNum(),
	}
}
