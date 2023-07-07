package src

import (
	"github.com/BurntSushi/toml"
)

const (
	PackageFile = "package.toml"
)

func parser() Package {
	var p Package
	if _, err := toml.DecodeFile(PackageFile, &p); err != nil {
		throw("package.toml parser failed: %v", err)
	}

	return p
}
