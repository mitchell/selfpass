package transport

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/mitchell/selfpass/credentials/endpoints"
	"github.com/mitchell/selfpass/credentials/protobuf"
	"github.com/mitchell/selfpass/credentials/types"
)

func decodeGetAllMetadataRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.GetAllMetadataRequest)
	return endpoints.GetAllMetadataRequest{
		SourceHost: r.SourceHost,
	}, nil
}

func encodeMetadataStreamResponse(ctx context.Context, response interface{}) (interface{}, error) {
	r := response.(endpoints.MetadataStream)
	pbmdch := make(chan protobuf.Metadata, 1)

	go func() {
		defer close(pbmdch)

		for md := range r.Metadata {
			createdAt, err := ptypes.TimestampProto(md.CreatedAt)
			if err != nil {
				r.Errors <- err
				return
			}

			updatedAt, err := ptypes.TimestampProto(md.UpdatedAt)
			if err != nil {
				r.Errors <- err
				return
			}

			pbmdch <- protobuf.Metadata{
				Id:         md.ID,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				SourceHost: md.SourceHost,
				Primary:    md.Primary,
				LoginUrl:   md.LoginURL,
			}
		}
	}()

	return protobufMetadataStream{
		Metadata: pbmdch,
		Errors:   r.Errors,
	}, nil
}

type protobufMetadataStream struct {
	Metadata <-chan protobuf.Metadata
	Errors   chan error
}

func decodeCredentialRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.CredentialRequest)
	return types.CredentialInput{
		MetadataInput: types.MetadataInput{
			Primary:    r.Primary,
			LoginURL:   r.LoginUrl,
			SourceHost: r.SourceHost,
		},
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}, nil
}

func encodeCredentialResponse(ctx context.Context, response interface{}) (interface{}, error) {
	r := response.(types.Credential)

	createdAt, err := ptypes.TimestampProto(r.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.TimestampProto(r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return protobuf.Credential{
		Id:         r.ID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Primary:    r.Primary,
		SourceHost: r.SourceHost,
		LoginUrl:   r.LoginURL,
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
	}, nil
}

func decodeUpdateRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.UpdateRequest)
	return endpoints.UpdateRequest{
		ID: r.Id,
		Credential: types.CredentialInput{
			MetadataInput: types.MetadataInput{
				Primary:    r.Credential.Primary,
				SourceHost: r.Credential.SourceHost,
				LoginURL:   r.Credential.LoginUrl,
			},
			Username: r.Credential.Username,
			Email:    r.Credential.Email,
			Password: r.Credential.Password,
		},
	}, nil
}

func decodeIdRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.IdRequest)
	return endpoints.IDRequest{
		ID: r.Id,
	}, nil
}

func noOpEncode(ctx context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}
