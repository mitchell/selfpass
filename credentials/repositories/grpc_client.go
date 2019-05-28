package repositories

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mitchell/selfpass/credentials/endpoints"
	"github.com/mitchell/selfpass/credentials/protobuf"
	"github.com/mitchell/selfpass/credentials/transport"
	"github.com/mitchell/selfpass/credentials/types"
)

func NewCredentialServiceClient(ctx context.Context, target, ca, cert, key string) (types.CredentialClient, error) {
	keypair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	capool := x509.NewCertPool()
	capool.AppendCertsFromPEM([]byte(ca))

	creds := credentials.NewTLS(&tls.Config{
		RootCAs:      capool,
		Certificates: []tls.Certificate{keypair},
	})

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return credentialServiceClient{
		client: protobuf.NewCredentialServiceClient(conn),
	}, nil
}

type credentialServiceClient struct {
	client protobuf.CredentialServiceClient
}

func (c credentialServiceClient) GetAllMetadata(ctx context.Context, sourceHost string) (output <-chan types.Metadata, errch chan error) {
	pbmdch := make(chan protobuf.Metadata, 1)
	errch = make(chan error, 1)

	stream, err := transport.DecodeMetdataStreamResponse(ctx, transport.ProtobufMetadataStream{
		Metadata: pbmdch,
		Errors:   errch,
	})

	srv, err := c.client.GetAllMetadata(ctx, &protobuf.GetAllMetadataRequest{SourceHost: sourceHost})
	if err != nil {
		errch <- err
		return nil, errch
	}

	go func() {
		defer close(pbmdch)

		for {
			select {
			case <-ctx.Done():
				errch <- fmt.Errorf("context timeout")
				return
			default:
			}

			pbmd, err := srv.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				errch <- err
				return
			}

			pbmdch <- *pbmd
		}
	}()

	return stream.Metadata, stream.Errors
}

func (c credentialServiceClient) Get(ctx context.Context, id string) (output types.Credential, err error) {
	req := transport.EncodeIdRequest(endpoints.IDRequest{ID: id})

	res, err := c.client.Get(ctx, &req)
	if err != nil {
		return output, err
	}

	return transport.DecodeCredential(*res)
}

func (c credentialServiceClient) Create(ctx context.Context, ci types.CredentialInput) (output types.Credential, err error) {
	req := transport.EncodeCredentialRequest(ci)

	res, err := c.client.Create(ctx, &req)
	if err != nil {
		return output, err
	}

	return transport.DecodeCredential(*res)
}

func (c credentialServiceClient) Update(ctx context.Context, id string, ci types.CredentialInput) (output types.Credential, err error) {
	req := transport.EncodeUpdateRequest(endpoints.UpdateRequest{ID: id, Credential: ci})

	res, err := c.client.Update(ctx, &req)
	if err != nil {
		return output, err
	}

	return transport.DecodeCredential(*res)
}

func (c credentialServiceClient) Delete(ctx context.Context, id string) (err error) {
	req := transport.EncodeIdRequest(endpoints.IDRequest{ID: id})

	res, err := c.client.Delete(ctx, &req)
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("delete unsuccessful")
	}
	return nil
}
