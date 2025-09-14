package grpcserver

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"user-service/internal/db"
	"user-service/internal/user"
	userv1 "user-service/pb/userv1"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
}

func New() *UserServiceServer { return &UserServiceServer{} }

func toPB(u *user.User) *userv1.User {
	return &userv1.User{
		Id:              u.ID.String(),
		Username:        u.Username,
		Name:            u.Name,
		Lastname:        u.Lastname,
		PhoneNumber:     u.PhoneNumber,
		Email:           u.Email,
		OauthProvider:   u.OAuthProvider,
		OauthProviderId: u.OAuthProviderID,
		CreatedAt:       timestamppb.New(u.CreatedAt),
		UpdatedAt:       timestamppb.New(u.UpdatedAt),
	}
}

func (s *UserServiceServer) UpsertUser(ctx context.Context, req *userv1.UpsertUserRequest) (*userv1.UpsertUserResponse, error) {
	// naive upsert by email
	var u user.User
	tx := db.Gorm()
	err := tx.Where("email = ?", req.Email).First(&u).Error
	if err != nil {
		// create
		u = user.User{
			Username:        req.GetUsername(),
			Name:            req.GetName(),
			Lastname:        req.GetLastname(),
			PhoneNumber:     req.GetPhoneNumber(),
			Email:           req.GetEmail(),
			OAuthProvider:   req.GetOauthProvider(),
			OAuthProviderID: req.GetOauthProviderId(),
		}
		if err := tx.Create(&u).Error; err != nil {
			return nil, err
		}
	} else {
		updates := map[string]any{}
		if v := req.GetUsername(); v != "" {
			updates["username"] = v
		}
		if v := req.GetName(); v != "" {
			updates["name"] = v
		}
		if v := req.GetLastname(); v != "" {
			updates["lastname"] = v
		}
		if v := req.GetPhoneNumber(); v != "" {
			updates["phone_number"] = v
		}
		if v := req.GetOauthProvider(); v != "" {
			updates["oauth_provider"] = v
		}
		if v := req.GetOauthProviderId(); v != "" {
			updates["oauth_provider_id"] = v
		}
		if len(updates) > 0 {
			if err := tx.Model(&u).Where("email = ?", req.Email).Updates(updates).Error; err != nil {
				return nil, err
			}
			// reload
			_ = tx.Where("email = ?", req.Email).First(&u).Error
		}
	}
	return &userv1.UpsertUserResponse{User: toPB(&u)}, nil
}

func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.GetUserByEmailResponse, error) {
	var u user.User
	if err := db.Gorm().Where("email = ?", req.GetEmail()).First(&u).Error; err != nil {
		return nil, err
	}
	return &userv1.GetUserByEmailResponse{User: toPB(&u)}, nil
}

func (s *UserServiceServer) UpdateUserProfile(ctx context.Context, req *userv1.UpdateUserProfileRequest) (*userv1.UpdateUserProfileResponse, error) {
	updates := map[string]any{}
	if v := req.GetUsername(); v != "" {
		updates["username"] = v
	}
	if v := req.GetName(); v != "" {
		updates["name"] = v
	}
	if v := req.GetLastname(); v != "" {
		updates["lastname"] = v
	}
	// allow clearing phone number: use a pointer or explicit empty allowed
	updates["phone_number"] = req.GetPhoneNumber()
	if len(updates) == 0 {
		// nothing to update, still return current record
		var u user.User
		if err := db.Gorm().Where("email = ?", req.GetEmail()).First(&u).Error; err != nil {
			return nil, err
		}
		return &userv1.UpdateUserProfileResponse{User: toPB(&u)}, nil
	}
	if err := db.Gorm().Model(&user.User{}).Where("email = ?", req.GetEmail()).Updates(updates).Error; err != nil {
		return nil, err
	}
	var u user.User
	if err := db.Gorm().Where("email = ?", req.GetEmail()).First(&u).Error; err != nil {
		return nil, err
	}
	// touch updatedAt to now if zero
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return &userv1.UpdateUserProfileResponse{User: toPB(&u)}, nil
}
