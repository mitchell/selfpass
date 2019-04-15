package types

import "time"

// Credential TODO
type Credential struct {
	Metadata
	Username string
	Email    string
	Password string
}

// CredentialInput TODO
type CredentialInput struct {
	MetadataInput
	Username string
	Email    string
	Password string
}

// Metadata TODO
type Metadata struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Primary    string
	SourceHost string
	LoginURL   string
}

// MetadataInput TODO
type MetadataInput struct {
	Primary    string
	SourceHost string
	LoginURL   string
}
