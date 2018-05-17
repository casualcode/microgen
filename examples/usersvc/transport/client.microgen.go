// Code generated by microgen 1.0.0-alpha. DO NOT EDIT.

package transport

import (
	"context"
	"errors"
	usersvc "github.com/devimteam/microgen/examples/usersvc/usersvc"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	opentracinggo "github.com/opentracing/opentracing-go"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// TraceClientEndpoints is used for tracing endpoints on client side.
func TraceClientEndpoints(endpoints EndpointsSet, tracer opentracinggo.Tracer) EndpointsSet {
	return EndpointsSet{
		CreateCommentEndpoint:   opentracing.TraceClient(tracer, "CreateComment")(endpoints.CreateCommentEndpoint),
		CreateUserEndpoint:      opentracing.TraceClient(tracer, "CreateUser")(endpoints.CreateUserEndpoint),
		FindUsersEndpoint:       opentracing.TraceClient(tracer, "FindUsers")(endpoints.FindUsersEndpoint),
		GetCommentEndpoint:      opentracing.TraceClient(tracer, "GetComment")(endpoints.GetCommentEndpoint),
		GetUserCommentsEndpoint: opentracing.TraceClient(tracer, "GetUserComments")(endpoints.GetUserCommentsEndpoint),
		GetUserEndpoint:         opentracing.TraceClient(tracer, "GetUser")(endpoints.GetUserEndpoint),
		UpdateUserEndpoint:      opentracing.TraceClient(tracer, "UpdateUser")(endpoints.UpdateUserEndpoint),
	}
}

func (set EndpointsSet) CreateUser(arg0 context.Context, arg1 usersvc.User) (res0 string, res1 error) {
	request := CreateUserRequest{User: arg1}
	response, res1 := set.CreateUserEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*CreateUserResponse).Id, res1
}

func (set EndpointsSet) UpdateUser(arg0 context.Context, arg1 usersvc.User) (res0 error) {
	request := UpdateUserRequest{User: arg1}
	_, res0 = set.UpdateUserEndpoint(arg0, &request)
	if res0 != nil {
		if e, ok := status.FromError(res0); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res0 = errors.New(e.Message())
		}
		return
	}
	return res0
}

func (set EndpointsSet) GetUser(arg0 context.Context, arg1 string) (res0 usersvc.User, res1 error) {
	request := GetUserRequest{Id: arg1}
	response, res1 := set.GetUserEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*GetUserResponse).User, res1
}

func (set EndpointsSet) FindUsers(arg0 context.Context) (res0 map[string]usersvc.User, res1 error) {
	request := FindUsersRequest{}
	response, res1 := set.FindUsersEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*FindUsersResponse).Results, res1
}

func (set EndpointsSet) CreateComment(arg0 context.Context, arg1 usersvc.Comment) (res0 string, res1 error) {
	request := CreateCommentRequest{Comment: arg1}
	response, res1 := set.CreateCommentEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*CreateCommentResponse).Id, res1
}

func (set EndpointsSet) GetComment(arg0 context.Context, arg1 string) (res0 usersvc.Comment, res1 error) {
	request := GetCommentRequest{Id: arg1}
	response, res1 := set.GetCommentEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*GetCommentResponse).Comment, res1
}

func (set EndpointsSet) GetUserComments(arg0 context.Context, arg1 string) (res0 []usersvc.Comment, res1 error) {
	request := GetUserCommentsRequest{UserId: arg1}
	response, res1 := set.GetUserCommentsEndpoint(arg0, &request)
	if res1 != nil {
		if e, ok := status.FromError(res1); ok || e.Code() == codes.Internal || e.Code() == codes.Unknown {
			res1 = errors.New(e.Message())
		}
		return
	}
	return response.(*GetUserCommentsResponse).List, res1
}
