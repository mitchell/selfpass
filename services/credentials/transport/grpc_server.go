package transport

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protobuf "github.com/mitchell/selfpass/protobuf/go"
	"github.com/mitchell/selfpass/services/credentials/endpoints"
	"github.com/mitchell/selfpass/services/credentials/types"
)

func NewGRPCServer(svc types.Service, logger log.Logger) GRPCServer {
	return GRPCServer{
		getAllMetadata: grpc.NewServer(
			endpoints.MakeGetAllMetadataEndpoint(svc),
			decodeSourceHostRequest,
			encodeMetadataStreamResponse,
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		get: grpc.NewServer(
			endpoints.MakeGetEndpoint(svc),
			decodeIdRequest,
			encodeCredentialResponse,
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		create: grpc.NewServer(
			endpoints.MakeCreateEndpoint(svc),
			decodeCredentialRequest,
			encodeCredentialResponse,
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		update: grpc.NewServer(
			endpoints.MakeUpdateEndpoint(svc),
			decodeUpdateRequest,
			encodeCredentialResponse,
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		delete: grpc.NewServer(
			endpoints.MakeDeleteEndpoint(svc),
			decodeIdRequest,
			noOp,
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
	}
}

type GRPCServer struct {
	getAllMetadata *grpc.Server
	get            *grpc.Server
	create         *grpc.Server
	update         *grpc.Server
	delete         *grpc.Server
}

func (s GRPCServer) GetAllMetadata(r *protobuf.SourceHostRequest, srv protobuf.Credentials_GetAllMetadataServer) (err error) {
	var i interface{}
	ctx := srv.Context()

	ctx, i, err = s.getAllMetadata.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return err
	}

	mds := i.(ProtobufMetadataStream)

receiveLoop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break receiveLoop
		case err = <-mds.Errors:
			if err != nil {
				err = handlerGRPCError(err)
				break receiveLoop
			}
		case md, ok := <-mds.Metadata:
			if !ok {
				break receiveLoop
			}
			if err = srv.Send(&md); err != nil {
				err = handlerGRPCError(err)
				break receiveLoop
			}
		}
	}

	return err
}

func (s GRPCServer) Get(ctx context.Context, r *protobuf.IdRequest) (*protobuf.Credential, error) {
	ctx, i, err := s.get.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	c := &protobuf.Credential{}
	*c = i.(protobuf.Credential)
	return c, nil
}

func (s GRPCServer) Create(ctx context.Context, r *protobuf.CredentialRequest) (*protobuf.Credential, error) {
	ctx, i, err := s.create.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	c := &protobuf.Credential{}
	*c = i.(protobuf.Credential)
	return c, nil
}

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

func (s GRPCServer) Delete(ctx context.Context, r *protobuf.IdRequest) (*protobuf.SuccessResponse, error) {
	ctx, _, err := s.delete.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	return &protobuf.SuccessResponse{Success: true}, nil
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
