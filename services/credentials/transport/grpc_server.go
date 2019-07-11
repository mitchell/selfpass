package transport

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchell/selfpass/services/credentials/endpoints"
	"github.com/mitchell/selfpass/services/credentials/protobuf"
	"github.com/mitchell/selfpass/services/credentials/types"
)

func NewGRPCServer(svc types.Service, logger log.Logger) GRPCServer {
	return GRPCServer{
		getAllMetadata: grpc.NewServer(
			endpoints.MakeGetAllMetadataEndpoint(svc),
			decodeGetAllMetadataRequest,
			encodeMetadataStreamResponse,
			grpc.ServerErrorLogger(logger),
		),
		get: grpc.NewServer(
			endpoints.MakeGetEndpoint(svc),
			decodeIdRequest,
			encodeCredentialResponse,
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
			noOp,
			grpc.ServerErrorLogger(logger),
		),
		dump: grpc.NewServer(
			endpoints.MakeDumpEndpoint(svc),
			noOp,
			encodeDumpResponse,
			grpc.ServerErrorLogger(logger),
		),
	}
}

type GRPCServer struct {
	getAllMetadata *grpc.Server
	get            *grpc.Server
	create         *grpc.Server
	update         *grpc.Server
	delete         *grpc.Server
	dump           *grpc.Server
}

func (s GRPCServer) GetAllMetadata(r *protobuf.GetAllMetadataRequest, srv protobuf.CredentialService_GetAllMetadataServer) (err error) {
	defer func() { err = handlerGRPCError(err) }()

	var i interface{}
	ctx := srv.Context()

	ctx, i, err = s.getAllMetadata.ServeGRPC(ctx, *r)
	if err != nil {
		return err
	}

	mds := i.(ProtobufMetadataStream)

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
			if err = srv.Send(&md); err != nil {
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

func (s GRPCServer) Delete(ctx context.Context, r *protobuf.IdRequest) (*protobuf.DeleteResponse, error) {
	ctx, _, err := s.delete.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	return &protobuf.DeleteResponse{Success: true}, nil
}

func (s GRPCServer) Dump(ctx context.Context, r *protobuf.EmptyRequest) (*protobuf.DumpResponse, error) {
	ctx, i, err := s.dump.ServeGRPC(ctx, *r)
	if err != nil {
		err = handlerGRPCError(err)
		return nil, err
	}

	res := &protobuf.DumpResponse{}
	*res = i.(protobuf.DumpResponse)
	return res, nil
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
