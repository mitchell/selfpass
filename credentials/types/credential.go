package types

import (
	"fmt"
	"time"
)

const TypePrefixCred = "cred"

type Credential struct {
	Metadata
	Username  string
	Email     string
	Password  string `json:"-"`
	OTPSecret string `json:"-"`
}

func (c Credential) String() string {
	return fmt.Sprintf(
		"username = %s\nemail = %s\n%s",
		c.Username, c.Email, c.Metadata,
	)
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
	return fmt.Sprintf(
		"id = %s\nsourceHost = %s\ncreatedAt = %s\nupdatedAt = %s\nprimary = %s\nloginUrl = %s\ntag = %s\n",
		m.ID, m.SourceHost, m.CreatedAt, m.UpdatedAt, m.Primary, m.LoginURL, m.Tag,
	)
}

type MetadataInput struct {
	Primary    string
	SourceHost string
	LoginURL   string
	Tag        string
}
