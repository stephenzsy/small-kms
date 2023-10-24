package base

import (
	"bytes"
	"encoding"
	"io"

	"github.com/google/uuid"
)

type storageLocator struct {
	NID uuid.UUID
	RID uuid.UUID
}

// MarshalText implements encoding.TextMarshaler.
func (s storageLocator) MarshalText() (text []byte, err error) {
	b := make([]byte, 0, 36+1+36)
	if bn, err := s.NID.MarshalText(); err != nil {
		return nil, err
	} else {
		b = append(b, bn...)
		b = append(b, ':')
	}
	if br, err := s.RID.MarshalText(); err != nil {
		return nil, err
	} else {
		b = append(b, br...)
	}
	return b, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (sl *storageLocator) UnmarshalText(text []byte) (err error) {
	b := bytes.SplitN(text, []byte(":"), 2)
	l := storageLocator{}
	if l.NID, err = uuid.ParseBytes(b[0]); err != nil {
		return err
	}
	if l.RID, err = uuid.ParseBytes(b[1]); err != nil {
		return err
	}
	*sl = l
	return nil
}

var _ encoding.TextMarshaler = storageLocator{}
var _ encoding.TextUnmarshaler = (*storageLocator)(nil)

func (sl *storageLocator) String() string {
	b, _ := sl.MarshalText()
	return string(b)
}

func (sl *storageLocator) IsNilOrEmpty() bool {
	if sl == nil {
		return true
	}
	return sl.NID == uuid.Nil && sl.RID == uuid.Nil
}

func (sl *storageLocator) WriteToDigest(writer io.Writer) (s int, err error) {
	if sl == nil {
		return 0, nil
	}
	if c, err := writer.Write(sl.NID[:]); err != nil {
		return s + c, err
	} else {
		s += c
	}
	if c, err := writer.Write(sl.RID[:]); err != nil {
		return s + c, err
	} else {
		s += c
	}
	return s, nil
}
