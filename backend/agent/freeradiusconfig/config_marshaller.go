package frconfig

import "strings"

type FreeRadiusConfigMarshaler interface {
	MarshalFreeradiusConfig(sb *strings.Builder, linePrefix string) error
}

type FreeRadiusConfigList[T FreeRadiusConfigMarshaler] []T

func (l FreeRadiusConfigList[T]) MarshalFreeradiusConfig(sb *strings.Builder, linePrefix string) error {
	for _, c := range l {
		err := c.MarshalFreeradiusConfig(sb, linePrefix)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ FreeRadiusConfigMarshaler = FreeRadiusConfigList[FreeRadiusConfigMarshaler](nil)
