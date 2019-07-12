package types

import "context"

type Service interface {
	GetAllMetadata(ctx context.Context, sourceHost string) (output <-chan Metadata, errch chan error)
	Get(ctx context.Context, id string) (output Credential, err error)
	Create(ctx context.Context, ci CredentialInput) (output Credential, err error)
	Update(ctx context.Context, id string, ci CredentialInput) (output Credential, err error)
	Delete(ctx context.Context, id string) (err error)
	DumpDB(ctx context.Context) (bs []byte, err error)
}

type CredentialsRepo interface {
	GetAllMetadata(ctx context.Context, sourceHost string, errch chan<- error) (output <-chan Metadata)
	Get(ctx context.Context, id string) (output Credential, err error)
	Put(ctx context.Context, c Credential) (err error)
	Delete(ctx context.Context, id string) (err error)
	DumpDB(ctx context.Context) (bs []byte, err error)
}

type CredentialsClientInit func(ctx context.Context, target, ca, cert, key string) (c CredentialsClient, err error)

type CredentialsClient interface {
	GetAllMetadata(ctx context.Context, sourceHost string) (output <-chan Metadata, errch chan error)
	Get(ctx context.Context, id string) (output Credential, err error)
	Create(ctx context.Context, ci CredentialInput) (output Credential, err error)
	Update(ctx context.Context, id string, ci CredentialInput) (output Credential, err error)
	Delete(ctx context.Context, id string) (err error)
}
