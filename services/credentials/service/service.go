package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/mitchell/selfpass/services/credentials/types"
)

func NewCredentials(repo types.CredentialsRepo) Credentials {
	return Credentials{
		repo: repo,
	}
}

type Credentials struct {
	repo types.CredentialsRepo
}

func (svc Credentials) GetAllMetadata(ctx context.Context, sourceHost string) (output <-chan types.Metadata, errch chan error) {
	errch = make(chan error, 1)
	output = svc.repo.GetAllMetadata(ctx, sourceHost, errch)
	return output, errch
}

func (svc Credentials) Get(ctx context.Context, id string) (output types.Credential, err error) {
	if id == "" {
		return output, fmt.Errorf("%s must specify an id", types.InvalidArgument)
	}
	return svc.repo.Get(ctx, id)
}

func (svc Credentials) Create(ctx context.Context, ci types.CredentialInput) (output types.Credential, err error) {
	if err = validateCredentialInput(ci); err != nil {
		return output, err
	}

	now := time.Now()

	var c types.Credential
	c.ID = generateID()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Primary = ci.Primary
	c.LoginURL = ci.LoginURL
	c.SourceHost = ci.SourceHost
	c.Username = ci.Username
	c.Email = ci.Email
	c.Password = ci.Password
	c.OTPSecret = ci.OTPSecret
	c.Tag = ci.Tag

	err = svc.repo.Put(ctx, c)

	return c, err
}

func validateCredentialInput(c types.CredentialInput) (err error) {
	switch {
	case c.SourceHost == "":
		return fmt.Errorf("%s must specify source host", types.InvalidArgument)
	case c.Primary == "":
		return fmt.Errorf("%s must specify primary user key", types.InvalidArgument)
	case c.Password == "":
		return fmt.Errorf("%s must specify password", types.InvalidArgument)
	}

	if _, err = base64.StdEncoding.DecodeString(c.Password); err != nil {
		return fmt.Errorf("%s password must be encrypted and base64 encoded", types.InvalidArgument)
	}

	if c.OTPSecret != "" {
		if _, err = base64.StdEncoding.DecodeString(c.OTPSecret); err != nil {
			return fmt.Errorf("%s otp secret must be encrypted and base64 encoded", types.InvalidArgument)
		}
	}

	return err
}

func generateID() string {
	const idLen = 8
	const alphanumerics = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ0123456789"
	const alphaLen = len(alphanumerics)

	rand.Seed(time.Now().UnixNano())
	id := make([]byte, idLen)

	for index := range id {
		id[index] = alphanumerics[rand.Int63()%int64(alphaLen)]
	}

	return fmt.Sprintf("%s-%s", types.KeyCredential, string(id))
}

func (svc Credentials) Update(ctx context.Context, id string, ci types.CredentialInput) (output types.Credential, err error) {
	if err = validateCredentialInput(ci); err != nil {
		return output, err
	}

	if id == "" {
		return output, fmt.Errorf("%s must specify an id", types.InvalidArgument)
	}

	c, err := svc.repo.Get(ctx, id)
	if err != nil {
		return output, err
	}

	c.UpdatedAt = time.Now()
	c.Primary = ci.Primary
	c.LoginURL = ci.LoginURL
	c.SourceHost = ci.SourceHost
	c.Password = ci.Password
	c.OTPSecret = ci.OTPSecret
	c.Email = ci.Email
	c.Username = ci.Username
	c.Tag = ci.Tag

	return c, svc.repo.Put(ctx, c)
}

func (svc Credentials) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		return fmt.Errorf("%s must specify an id", types.InvalidArgument)
	}
	return svc.repo.Delete(ctx, id)
}
