package omaha

import (
	"encoding/base64"
	"encoding/xml"
)

type Hash []byte

func (h *Hash) UnmarshalXMLAttr(attr xml.Attr) error {
	b, err := base64.StdEncoding.DecodeString(attr.Value)
	if err != nil {
		return err
	}
	*h = b
	return nil
}

func (h Hash) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: base64.StdEncoding.EncodeToString(h),
	}, nil
}
