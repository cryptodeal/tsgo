package tsgo

import "strings"

func (g *PackageGenerator) IsEnumStruct(name string) bool {
	for k, v := range g.conf.EnumStructs {
		if strings.EqualFold(k, name) || strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}

func (g *PackageGenerator) LastFFIFunc() string {
	var lastKey string
	for k := range g.ffi.FFIFuncs {
		lastKey = k
	}
	return lastKey
}
