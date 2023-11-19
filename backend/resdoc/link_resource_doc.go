package resdoc

import "github.com/stephenzsy/small-kms/backend/models"

type LinkResourceDoc struct {
	ResourceDoc
	LinkTo DocIdentifier `json:"linkTo"`
}

func (doc *LinkResourceDoc) ToModel() (m models.LinkRef) {
	m.Ref = doc.ResourceDoc.ToRef()
	m.LinkTo = doc.LinkTo.String()
	return m
}
