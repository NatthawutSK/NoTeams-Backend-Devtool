package usersUsecase

import (
	"fmt"
	"mime/multipart"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersRepository"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/auth"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	DeleteOauth(oauthId string) (string, error)
	GetUserProfile(userId string) (*users.User, error)
	InsertUser(req *users.UserRegisterReq) (*users.User, error)
	GetPassport(req *users.UserLoginReq) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredentialReq) (*users.UserPassport, error)
	FindByEmailOrUsername(email, username string) (*users.FindMember, error)
	UpdateUserProfile(userId string, req *users.UserUpdateProfileReq, avatarFile []*multipart.FileHeader) (*users.User, error)
	GetTeamsByUserId(userId string) ([]*users.TeamsByUserIdRes, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepository.IUserRepository
	upload          utils.IUpload
}

func UserUsecase(usersRepo usersRepository.IUserRepository, cfg config.IConfig) IUserUsecase {
	return &usersUsecase{
		usersRepository: usersRepo,
		cfg:             cfg,
		upload:          utils.Upload(cfg),
	}
}

// use for register user
func (u *usersUsecase) InsertUser(req *users.UserRegisterReq) (*users.User, error) {
	//hashing password
	if err := utils.BcryptHashing(req); err != nil {
		return nil, err
	}
	//insert user
	result, err := u.usersRepository.InsertUser(req)
	if err != nil {
		return nil, err
	}
	res, err := result.Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// use for login to get token and user information
func (u *usersUsecase) GetPassport(req *users.UserLoginReq) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmailOrUsername(req.Email, "")
	if err != nil {
		return nil, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// sign token
	accessToken, err1 := auth.NewRiAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.UserId,
	})
	if err1 != nil {
		return nil, err
	}
	refreshToken, err2 := auth.NewRiAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.UserId,
	})
	if err2 != nil {
		return nil, err
	}

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			UserId:   user.UserId,
			Email:    user.Email,
			Username: user.Username,
			Dob:      user.Dob,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	fmt.Println(passport.Token.RefreshToken)

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil

}

// use for logout
func (u *usersUsecase) DeleteOauth(oauthId string) (string, error) {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return "", err
	}
	return "logout success", nil

}

// use for refresh token
func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredentialReq) (*users.UserPassport, error) {
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id: profile.UserId,
	}

	accessToken, err := auth.NewRiAuth(auth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}
	refreshToken := auth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			OauthId:      oauth.OauthId,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil

}

// use for get user profile
func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil

}

func (u *usersUsecase) FindByEmailOrUsername(email, username string) (*users.FindMember, error) {
	member, err := u.usersRepository.FindOneUserByEmailOrUsername(email, username)
	if err != nil {
		return nil, err
	}

	res := &users.FindMember{
		UserId:   member.UserId,
		Username: member.Username,
		Avatar:   member.Avatar,
		Email:    member.Email,
	}

	return res, nil
}

func (u *usersUsecase) UpdateUserProfile(userId string, req *users.UserUpdateProfileReq, avatarFile []*multipart.FileHeader) (*users.User, error) {
	//check avatarFile
	if len(avatarFile) > 0 {
		//check len of avatarFile must be 1
		if len(avatarFile) > 1 {
			return nil, fmt.Errorf("avatar file must be 1")
		}

		//upload avatar
		url, err := u.upload.UploadFiles(avatarFile, false, userId)
		if err != nil {
			return nil, err
		}

		req.AvatarUrl = url[0].Url
	}

	//update profile
	if err := u.usersRepository.UpdateUserProfile(userId, req); err != nil {
		return nil, err
	}

	//get profile
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *usersUsecase) GetTeamsByUserId(userId string) ([]*users.TeamsByUserIdRes, error) {
	teams, err := u.usersRepository.GetTeamsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return teams, nil
}
