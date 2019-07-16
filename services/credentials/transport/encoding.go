package transport

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	protobuf "github.com/mitchell/selfpass/protobuf/go"
	"github.com/mitchell/selfpass/services/credentials/endpoints"
	"github.com/mitchell/selfpass/services/credentials/types"
)

func decodeSourceHostRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.SourceHostRequest)
	return endpoints.SourceHostRequest{
		SourceHost: r.SourceHost,
	}, nil
}

func EncodeSourceHostRequest(request endpoints.SourceHostRequest) *protobuf.SourceHostRequest {
	return &protobuf.SourceHostRequest{
		SourceHost: request.SourceHost,
	}
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
				Tag:        md.Tag,
			}
		}
	}()

	return ProtobufMetadataStream{
		Metadata: pbmdch,
		Errors:   r.Errors,
	}, nil
}

func DecodeMetdataStreamResponse(ctx context.Context, r ProtobufMetadataStream) (endpoints.MetadataStream, error) {
	mdch := make(chan types.Metadata, 1)
	errch := make(chan error, 1)

	go func() {
		defer close(mdch)

		for pbmd := range r.Metadata {
			createdAt, err := ptypes.Timestamp(pbmd.CreatedAt)
			if err != nil {
				errch <- err
				return
			}

			updatedAt, err := ptypes.Timestamp(pbmd.UpdatedAt)
			if err != nil {
				errch <- err
				return
			}

			mdch <- types.Metadata{
				ID:         pbmd.Id,
				SourceHost: pbmd.SourceHost,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Primary:    pbmd.Primary,
				LoginURL:   pbmd.LoginUrl,
				Tag:        pbmd.Tag,
			}
		}
	}()

	return endpoints.MetadataStream{
		Metadata: mdch,
		Errors:   errch,
	}, nil
}

type ProtobufMetadataStream struct {
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
			Tag:        r.Tag,
		},
		Username:  r.Username,
		Email:     r.Email,
		Password:  r.Password,
		OTPSecret: r.OtpSecret,
	}, nil
}

func EncodeCredentialRequest(r types.CredentialInput) protobuf.CredentialRequest {
	return protobuf.CredentialRequest{
		Primary:    r.Primary,
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		OtpSecret:  r.OTPSecret,
		SourceHost: r.SourceHost,
		LoginUrl:   r.LoginURL,
		Tag:        r.Tag,
	}
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
		Tag:        r.Tag,
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		OtpSecret:  r.OTPSecret,
	}, nil
}

func DecodeCredential(r protobuf.Credential) (c types.Credential, err error) {

	createdAt, err := ptypes.Timestamp(r.CreatedAt)
	if err != nil {
		return c, err
	}

	updatedAt, err := ptypes.Timestamp(r.UpdatedAt)
	if err != nil {
		return c, err
	}

	return types.Credential{
		Metadata: types.Metadata{
			ID:         r.Id,
			SourceHost: r.SourceHost,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			Primary:    r.Primary,
			LoginURL:   r.LoginUrl,
			Tag:        r.Tag,
		},
		Username:  r.Username,
		Email:     r.Email,
		Password:  r.Password,
		OTPSecret: r.OtpSecret,
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
				Tag:        r.Credential.Tag,
			},
			Username:  r.Credential.Username,
			Email:     r.Credential.Email,
			Password:  r.Credential.Password,
			OTPSecret: r.Credential.OtpSecret,
		},
	}, nil
}

func EncodeUpdateRequest(r endpoints.UpdateRequest) protobuf.UpdateRequest {
	c := r.Credential
	return protobuf.UpdateRequest{
		Id: r.ID,
		Credential: &protobuf.CredentialRequest{
			Primary:    c.Primary,
			Username:   c.Username,
			Email:      c.Email,
			Password:   c.Password,
			OtpSecret:  c.OTPSecret,
			SourceHost: c.SourceHost,
			LoginUrl:   c.LoginURL,
			Tag:        c.Tag,
		},
	}
}

func decodeIdRequest(ctx context.Context, request interface{}) (interface{}, error) {
	r := request.(protobuf.IdRequest)
	return endpoints.IDRequest{
		ID: r.Id,
	}, nil
}

func EncodeIdRequest(r endpoints.IDRequest) protobuf.IdRequest {
	return protobuf.IdRequest{
		Id: r.ID,
	}
}

func noOp(context.Context, interface{}) (interface{}, error) {
	return nil, nil
}
