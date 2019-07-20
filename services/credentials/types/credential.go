package types

import (
	"fmt"
	"time"
)

const KeyCredential = "cred"

type Credential struct {
	Metadata
	Username  string
	Email     string
	Password  string `json:"-"`
	OTPSecret string `json:"-"`
}

func (c Credential) String() string {
	format := "%s"
	args := []interface{}{c.Metadata}

	if c.Username != "" {
		format += "username = %s\n"
		args = append(args, c.Username)
	}

	if c.Email != "" {
		format += "email = %s\n"
		args = append(args, c.Email)
	}

	return fmt.Sprintf(format, args...)
}

type CredentialInput struct {
	MetadataInput
	Username  string
	Email     string
	Password  string
	OTPSecret string
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

func (m Metadata) String() string {
	format := "id = %s\ncreatedAt = %s\nupdatedAt = %s\nsourceHost = %s\nprimary = %s\n"
	args := []interface{}{m.ID, m.CreatedAt, m.UpdatedAt, m.SourceHost, m.Primary}

	if m.LoginURL != "" {
		format += "loginUrl = %s\n"
		args = append(args, m.LoginURL)
	}

	if m.Tag != "" {
		format += "tag = %s\n"
		args = append(args, m.Tag)
	}

	return fmt.Sprintf(format, args...)
}

type MetadataInput struct {
	Primary    string
	SourceHost string
	LoginURL   string
	Tag        string
}
