package serializers

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user"
	pb "github.com/kitanoyoru/kitaDriveBot/protos/gen/go/user/v1"
)

func UserToProto(u user.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func UserFromProto(u *pb.User) user.User {
	return user.User{
		ID:        u.Id,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}
