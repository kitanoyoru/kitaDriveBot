package user

import (
	"context"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/serializers"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user"
	pb "github.com/kitanoyoru/kitaDriveBot/protos/gen/go/user/v1"
)

func NewUserServiceServer(service user.Service) pb.UserServiceServer {
	return &userServiceServer{
		service: service,
	}
}

type userServiceServer struct {
	pb.UnimplementedUserServiceServer

	service user.Service
}

func (s userServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersCall_Request) (*pb.ListUsersCall_Response, error) {
	users, err := s.service.ListUsers(ctx, user.ListUsersRequest{
		IDs: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, item := range users {
		pbUsers = append(pbUsers, serializers.UserToProto(item))
	}

	return &pb.ListUsersCall_Response{
		Users: pbUsers,
	}, nil
}

func (s userServiceServer) GetUser(ctx context.Context, req *pb.GetUserCall_Request) (*pb.GetUserCall_Response, error) {
	u, err := s.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserCall_Response{
		User: serializers.UserToProto(u),
	}, nil
}

func (s userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserCall_Request) (*pb.CreateUserCall_Response, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	u, err := s.service.CreateUser(ctx, user.CreateUserRequest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserCall_Response{
		User: serializers.UserToProto(u),
	}, nil
}

func (s userServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserCall_Request) (*pb.UpdateUserCall_Response, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	u, err := s.service.UpdateUser(ctx, user.UpdateUserRequest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserCall_Response{
		User: serializers.UserToProto(u),
	}, nil
}

func (s userServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserCall_Request) (*pb.DeleteUserCall_Response, error) {
	err := s.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserCall_Response{}, nil
}
