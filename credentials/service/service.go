package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mitchell/selfpass/credentials/types"
)

// NewCredentials TODO
func NewCredentials(repo types.CredentialRepo) Credentials {
	return Credentials{
		repo: repo,
	}
}

// Credentials TODO
type Credentials struct {
	repo types.CredentialRepo
}

// GetAllMetadata TODO
func (svc Credentials) GetAllMetadata(ctx context.Context, sourceService string) (output <-chan types.Metadata, errch chan error) {
	errch = make(chan error, 1)
	output = svc.repo.GetAllMetadata(ctx, sourceService, errch)
	return output, errch
}

// Get TODO
func (svc Credentials) Get(ctx context.Context, id string) (output types.Credential, err error) {
	return svc.repo.Get(nil, id)
}

// Create TODO
func (svc Credentials) Create(ctx context.Context, ci types.CredentialInput) (output types.Credential, err error) {
	if err = validateCreate(ci); err != nil {
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

func validateCreate(c types.CredentialInput) (err error) {
	return err
}

// Update TODO
func (svc Credentials) Update(ctx context.Context, id string, ci types.CredentialInput) (output types.Credential, err error) {
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

// Delete TODO
func (svc Credentials) Delete(ctx context.Context, id string) (err error) {
	return svc.repo.Delete(nil, id)
}
