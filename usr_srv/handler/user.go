package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/usr_srv/global"
	"mxshop_srvs/usr_srv/model"
	"mxshop_srvs/usr_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user *model.User) *proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		//BirthDay: user.BirthDay,
		Gender: user.Gender,
		Role:   int32(user.Role),
	}
	if user.BirthDay != nil {
		userInfoRsp.BirthDay = uint64(user.BirthDay.Unix())
	}
	return &userInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u UserServer) GetUserList(ctx context.Context, info *proto.PageInfo) (*proto.UserListResponse, error) {
	zap.S().Info("GetUserList()")
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(info.Pn), int(info.PSize))).Find(&users)
	for _, user := range users {
		userInfoResponse := ModelToResponse(&user)
		rsp.Data = append(rsp.Data, userInfoResponse)
	}
	return &rsp, nil
}

func (u UserServer) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 通过手机号码查询用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(&user)
	return userInfoRsp, nil
}

func (u UserServer) GerUserById(ctx context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(&user)
	return userInfoRsp, nil
}

func (u UserServer) CreateUser(ctx context.Context, info *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: info.Mobile}).First(&user)
	if result.RowsAffected >= 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = info.Mobile
	user.NickName = info.NickName

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(info.PassWord, options)
	user.Password = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(&user), nil
}

func (u UserServer) UpdateUser(ctx context.Context, info *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, info.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthDay := time.Unix(int64(info.BirthDay), 0)
	user.NickName = info.NickName
	user.BirthDay = &birthDay
	user.Gender = info.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u UserServer) CheckPassword(ctx context.Context, info *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	passwordInfo := strings.Split(info.EncrpytedPassword, "$")
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	check := password.Verify(info.Password, passwordInfo[1], passwordInfo[2], options)
	return &proto.CheckResponse{Success: check}, nil
}
