package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mitchell/selfpass/credentials/types"
)

// MakeGetAllMetadataEndpoint TODO
func MakeGetAllMetadataEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(GetAllMetadataRequest)

		mdch, errch := svc.GetAllMetadata(ctx, r.SourceHost)

		return MetadataStream{
			Metadata: mdch,
			Errors:   errch,
		}, nil
	}
}

// MakeCreateEndpoint TODO
func MakeCreateEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(types.CredentialInput)
		return svc.Create(ctx, r)
	}
}

// MakeUpdateEndpoint TODO
func MakeUpdateEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UpdateRequest)
		return svc.Update(ctx, r.ID, r.Credential)
	}
}

func MakeDeleteEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(IDRequest)
		return nil, svc.Delete(ctx, r.ID)
	}
}

// IDRequest TODO
type IDRequest struct {
	ID string
}

// GetAllMetadataRequest TODO
type GetAllMetadataRequest struct {
	SourceHost string
}

// MetadataStream TODO
type MetadataStream struct {
	Metadata <-chan types.Metadata
	Errors   chan error
}

// UpdateRequest TODO
type UpdateRequest struct {
	ID         string
	Credential types.CredentialInput
}
