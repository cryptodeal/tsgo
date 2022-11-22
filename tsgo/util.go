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

func (g *PackageGenerator) IsWrappedEnum(name string) bool {
	for _, v := range g.conf.WrapStructs {
		if strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}
