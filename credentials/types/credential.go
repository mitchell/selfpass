package types

import (
	"time"
)

type Credential struct {
	Metadata
	Username string
	Email    string
	Password string
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
}

type MetadataInput struct {
	Primary    string
	SourceHost string
	LoginURL   string
}
