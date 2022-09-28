// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: tavern/v1/tavern.proto

package tavernv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/lnsp/ftp2p/gen/tavern/v1"
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
	// TavernServiceName is the fully-qualified name of the TavernService service.
	TavernServiceName = "tavern.v1.TavernService"
)

// TavernServiceClient is a client for the tavern.v1.TavernService service.
type TavernServiceClient interface {
	List(context.Context, *connect_go.Request[v1.ListRequest]) (*connect_go.Response[v1.ListResponse], error)
	Announce(context.Context, *connect_go.Request[v1.AnnounceRequest]) (*connect_go.Response[v1.AnnounceResponse], error)
}

// NewTavernServiceClient constructs a client for the tavern.v1.TavernService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTavernServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) TavernServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &tavernServiceClient{
		list: connect_go.NewClient[v1.ListRequest, v1.ListResponse](
			httpClient,
			baseURL+"/tavern.v1.TavernService/List",
			opts...,
		),
		announce: connect_go.NewClient[v1.AnnounceRequest, v1.AnnounceResponse](
			httpClient,
			baseURL+"/tavern.v1.TavernService/Announce",
			opts...,
		),
	}
}

// tavernServiceClient implements TavernServiceClient.
type tavernServiceClient struct {
	list     *connect_go.Client[v1.ListRequest, v1.ListResponse]
	announce *connect_go.Client[v1.AnnounceRequest, v1.AnnounceResponse]
}

// List calls tavern.v1.TavernService.List.
func (c *tavernServiceClient) List(ctx context.Context, req *connect_go.Request[v1.ListRequest]) (*connect_go.Response[v1.ListResponse], error) {
	return c.list.CallUnary(ctx, req)
}

// Announce calls tavern.v1.TavernService.Announce.
func (c *tavernServiceClient) Announce(ctx context.Context, req *connect_go.Request[v1.AnnounceRequest]) (*connect_go.Response[v1.AnnounceResponse], error) {
	return c.announce.CallUnary(ctx, req)
}

// TavernServiceHandler is an implementation of the tavern.v1.TavernService service.
type TavernServiceHandler interface {
	List(context.Context, *connect_go.Request[v1.ListRequest]) (*connect_go.Response[v1.ListResponse], error)
	Announce(context.Context, *connect_go.Request[v1.AnnounceRequest]) (*connect_go.Response[v1.AnnounceResponse], error)
}

// NewTavernServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTavernServiceHandler(svc TavernServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/tavern.v1.TavernService/List", connect_go.NewUnaryHandler(
		"/tavern.v1.TavernService/List",
		svc.List,
		opts...,
	))
	mux.Handle("/tavern.v1.TavernService/Announce", connect_go.NewUnaryHandler(
		"/tavern.v1.TavernService/Announce",
		svc.Announce,
		opts...,
	))
	return "/tavern.v1.TavernService/", mux
}

// UnimplementedTavernServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTavernServiceHandler struct{}

func (UnimplementedTavernServiceHandler) List(context.Context, *connect_go.Request[v1.ListRequest]) (*connect_go.Response[v1.ListResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("tavern.v1.TavernService.List is not implemented"))
}

func (UnimplementedTavernServiceHandler) Announce(context.Context, *connect_go.Request[v1.AnnounceRequest]) (*connect_go.Response[v1.AnnounceResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("tavern.v1.TavernService.Announce is not implemented"))
}