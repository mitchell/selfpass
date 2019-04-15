package types

import "context"

// Service TODO
type Service interface {
	GetAllMetadata(ctx context.Context, sourceService string) (output <-chan Metadata, errch chan error)
	Get(ctx context.Context, id string) (output Credential, err error)
	Create(ctx context.Context, ci CredentialInput) (output Credential, err error)
	Update(ctx context.Context, id string, ci CredentialInput) (output Credential, err error)
	Delete(ctx context.Context, id string) (err error)
}

// CredentialRepo TODO
type CredentialRepo interface {
	GetAllMetadata(ctx context.Context, sourceService string, errch chan<- error) (output <-chan Metadata)
	Get(ctx context.Context, id string) (output Credential, err error)
	Put(ctx context.Context, c Credential) (err error)
	Delete(ctx context.Context, id string) (err error)
}
