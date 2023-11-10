package frconfig

import "strings"

type FreeRadiusConfigMarshaler interface {
	MarshalFreeradiusConfig() ([]byte, error)
}

type FreeRadiusConfigList[T FreeRadiusConfigMarshaler] []T

func (l FreeRadiusConfigList[T]) MarshalFreeradiusConfig() ([]byte, error) {
	sb := &strings.Builder{}
	for i, c := range l {
		b, err := c.MarshalFreeradiusConfig()
		if err != nil {
			return nil, err
		}
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.Write(b)
	}
	return []byte(sb.String()), nil
}

var _ FreeRadiusConfigMarshaler = FreeRadiusConfigList[FreeRadiusConfigMarshaler](nil)
