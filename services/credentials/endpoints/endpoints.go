package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mitchell/selfpass/services/credentials/types"
)

func MakeCreateEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(types.CredentialInput)
		return svc.Create(ctx, r)
	}
}

func MakeDeleteEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(IDRequest)
		return nil, svc.Delete(ctx, r.ID)
	}
}

func MakeGetEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(IDRequest)
		return svc.Get(ctx, r.ID)
	}
}

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

func MakeUpdateEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UpdateRequest)
		return svc.Update(ctx, r.ID, r.Credential)
	}
}

func MakeDumpEndpoint(svc types.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		contents, err := svc.DumpDB(ctx)
		return DumpResponse{Contents: contents}, err
	}
}

type DumpResponse struct {
	Contents []byte
}

type IDRequest struct {
	ID string
}

type GetAllMetadataRequest struct {
	SourceHost string
}

type MetadataStream struct {
	Metadata <-chan types.Metadata
	Errors   chan error
}

type UpdateRequest struct {
	ID         string
	Credential types.CredentialInput
}
