package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mitchell/selfpass/credentials/types"
)

func NewCredentials(repo types.CredentialRepo) Credentials {
	return Credentials{
		repo: repo,
	}
}

type Credentials struct {
	repo types.CredentialRepo
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

	var c types.Credential
	c.ID = "cred-" + uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Primary = ci.Primary
	c.LoginURL = ci.LoginURL
	c.SourceHost = ci.SourceHost
	c.Username = ci.Username
	c.Email = ci.Email
	c.Password = ci.Password

	err = svc.repo.Put(ctx, c)

	return c, err
}

func validateCredentialInput(c types.CredentialInput) (err error) {
	switch {
	case c.SourceHost == "":
		return fmt.Errorf("%s must specify source host", types.InvalidArgument)
	case c.Password == "":
		return fmt.Errorf("%s must specify password", types.InvalidArgument)
	}

	return err
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
	c.Email = ci.Email
	c.Username = ci.Username

	return c, svc.repo.Put(ctx, c)
}

func (svc Credentials) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		return fmt.Errorf("%s must specify an id", types.InvalidArgument)
	}
	return svc.repo.Delete(ctx, id)
}

func (svc Credentials) DumpDB(ctx context.Context) (bs []byte, err error) {
	return svc.repo.DumpDB(ctx)
}
