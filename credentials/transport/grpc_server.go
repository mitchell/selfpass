package transport

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchell/selfpass/credentials/endpoints"
	"github.com/mitchell/selfpass/credentials/protobuf"
	"github.com/mitchell/selfpass/credentials/types"
)

// NewGRPCServer TODO
func NewGRPCServer(svc types.Service, logger log.Logger) GRPCServer {
	return GRPCServer{
		getAllMetadata: grpc.NewServer(
			endpoints.MakeGetAllMetadataEndpoint(svc),
			decodeGetAllMetadataRequest,
			encodeMetadataStreamResponse,
			grpc.ServerErrorLogger(logger),
		),
		create: grpc.NewServer(
			endpoints.MakeCreateEndpoint(svc),
			decodeCredentialRequest,
			encodeCredentialResponse,
			grpc.ServerErrorLogger(logger),
		),
		update: grpc.NewServer(
			endpoints.MakeUpdateEndpoint(svc),
			decodeUpdateRequest,
			encodeCredentialResponse,
			grpc.ServerErrorLogger(logger),
		),
		delete: grpc.NewServer(
			endpoints.MakeDeleteEndpoint(svc),
			decodeIdRequest,
			noOpEncode,
			grpc.ServerErrorLogger(logger),
		),
	}
}

// GRPCServer TODO
type GRPCServer struct {
	getAllMetadata *grpc.Server
	create         *grpc.Server
	update         *grpc.Server
	delete         *grpc.Server
}

// GetAllMetadata TODO
func (s GRPCServer) GetAllMetadata(r *protobuf.GetAllMetadataRequest, srv protobuf.CredentialService_GetAllMetadataServer) (err error) {
	defer func() {
		err = handlerGRPCError(err)
	}()

	var i interface{}
	ctx := srv.Context()

	ctx, i, err = s.getAllMetadata.ServeGRPC(ctx, *r)
	if err != nil {
		return err
	}

	mds := i.(protobufMetadataStream)

receiveLoop:
	for {
		select {
		case <-ctx.Done():
			break receiveLoop
		case err = <-mds.Errors:
			break receiveLoop
		case md, ok := <-mds.Metadata:
			if !ok {
				break receiveLoop
			}
			fmt.Println(md)
			if err = srv.Send(&md); err != nil {
				break receiveLoop
			}
		}
	}

	return err
}

// Get TODO
func (s GRPCServer) Get(context.Context, *protobuf.IdRequest) (*protobuf.Credential, error) {
	panic("implement me")
}

// Create TODO
func (s GRPCServer) Create(ctx context.Context, ci *protobuf.CredentialRequest) (*protobuf.Credential, error) {
	ctx, i, err := s.create.ServeGRPC(ctx, *ci)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	c := &protobuf.Credential{}
	*c = i.(protobuf.Credential)
	return c, nil
}

// Update TODO
func (s GRPCServer) Update(ctx context.Context, r *protobuf.UpdateRequest) (*protobuf.Credential, error) {
	ctx, i, err := s.update.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	c := &protobuf.Credential{}
	*c = i.(protobuf.Credential)
	return c, nil
}

// Delete TODO
func (s GRPCServer) Delete(ctx context.Context, r *protobuf.IdRequest) (*protobuf.DeleteResponse, error) {
	ctx, _, err := s.delete.ServeGRPC(ctx, *r)
	if err != nil {
		return nil, err
	}

	return &protobuf.DeleteResponse{Success: true}, nil
}

func handlerGRPCError(err error) error {
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), types.InvalidArgument):
			err = status.Error(codes.InvalidArgument, err.Error())
		case strings.HasPrefix(err.Error(), types.NotFound):
			err = status.Error(codes.NotFound, err.Error())
		default:
			err = status.Error(codes.Internal, "an internal error has occurred")
		}
	}
	return err
}
