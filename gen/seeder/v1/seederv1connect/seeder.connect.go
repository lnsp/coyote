// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: seeder/v1/seeder.proto

package seederv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/lnsp/ftp2p/gen/seeder/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// SeederServiceName is the fully-qualified name of the SeederService service.
	SeederServiceName = "seeder.v1.SeederService"
)

// SeederServiceClient is a client for the seeder.v1.SeederService service.
type SeederServiceClient interface {
	Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error)
	Fetch(context.Context, *connect_go.Request[v1.FetchRequest]) (*connect_go.Response[v1.FetchResponse], error)
}

// NewSeederServiceClient constructs a client for the seeder.v1.SeederService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewSeederServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) SeederServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &seederServiceClient{
		has: connect_go.NewClient[v1.HasRequest, v1.HasResponse](
			httpClient,
			baseURL+"/seeder.v1.SeederService/Has",
			opts...,
		),
		fetch: connect_go.NewClient[v1.FetchRequest, v1.FetchResponse](
			httpClient,
			baseURL+"/seeder.v1.SeederService/Fetch",
			opts...,
		),
	}
}

// seederServiceClient implements SeederServiceClient.
type seederServiceClient struct {
	has   *connect_go.Client[v1.HasRequest, v1.HasResponse]
	fetch *connect_go.Client[v1.FetchRequest, v1.FetchResponse]
}

// Has calls seeder.v1.SeederService.Has.
func (c *seederServiceClient) Has(ctx context.Context, req *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error) {
	return c.has.CallUnary(ctx, req)
}

// Fetch calls seeder.v1.SeederService.Fetch.
func (c *seederServiceClient) Fetch(ctx context.Context, req *connect_go.Request[v1.FetchRequest]) (*connect_go.Response[v1.FetchResponse], error) {
	return c.fetch.CallUnary(ctx, req)
}

// SeederServiceHandler is an implementation of the seeder.v1.SeederService service.
type SeederServiceHandler interface {
	Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error)
	Fetch(context.Context, *connect_go.Request[v1.FetchRequest]) (*connect_go.Response[v1.FetchResponse], error)
}

// NewSeederServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewSeederServiceHandler(svc SeederServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/seeder.v1.SeederService/Has", connect_go.NewUnaryHandler(
		"/seeder.v1.SeederService/Has",
		svc.Has,
		opts...,
	))
	mux.Handle("/seeder.v1.SeederService/Fetch", connect_go.NewUnaryHandler(
		"/seeder.v1.SeederService/Fetch",
		svc.Fetch,
		opts...,
	))
	return "/seeder.v1.SeederService/", mux
}

// UnimplementedSeederServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedSeederServiceHandler struct{}

func (UnimplementedSeederServiceHandler) Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("seeder.v1.SeederService.Has is not implemented"))
}

func (UnimplementedSeederServiceHandler) Fetch(context.Context, *connect_go.Request[v1.FetchRequest]) (*connect_go.Response[v1.FetchResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("seeder.v1.SeederService.Fetch is not implemented"))
}
