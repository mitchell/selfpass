package types

import (
	"time"
)

const TypePrefixCred = "cred"

type Credential struct {
	Metadata
	Username string
	Email    string
	Password string `json:"-"`
}

type CredentialInput struct {
	MetadataInput
	Username string
	Email    string
	Password string
}

type Metadata struct {
	ID         string // primary key
	SourceHost string // sort key
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Primary    string
	LoginURL   string
	Tag        string
}

type MetadataInput struct {
	Primary    string
	SourceHost string
	LoginURL   string
	Tag        string
}
