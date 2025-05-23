package service

import (
	"back/models"
	"back/models/dto"
	"back/repository"
)

type FriendService struct {
	authRepository   *repository.AuthRepository
	friendRepository *repository.FriendRepository
}

func NewFriendService(authRepository *repository.AuthRepository, friendRepository *repository.FriendRepository) *FriendService {
	return &FriendService{authRepository: authRepository, friendRepository: friendRepository}
}

func (s *FriendService) GetFriendInfo(friendId int) (dto.UserDto, error) {
	user, err := s.authRepository.SelectUserById(friendId)
	if user == (models.User{}) {
		return dto.UserDto{}, err
	}
	userDto := dto.UserDto{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Description: user.Description,
		Pic:         user.Pic,
	}
	return userDto, nil
}

func (s *FriendService) GetFriendList(userId int) ([]dto.UserDto, error) {
	friendList, err := s.friendRepository.GetFriendList(userId)
	if err != nil {
		return nil, err
	}
	var friends []dto.UserDto
	for _, friend := range friendList {
		userDto := dto.UserDto{
			Id:          friend.Id,
			Name:        friend.Name,
			Email:       friend.Email,
			Description: friend.Description,
			Pic:         friend.Pic,
		}
		friends = append(friends, userDto)
	}

	return friends, nil
}

func (s *FriendService) SearchUserByEmail(email string) (dto.UserDto, error) {
	user, err := s.authRepository.SelectUserByEmail(email)
	if user == (models.User{}) || err != nil {
		return dto.UserDto{}, err
	}
	userDto := dto.UserDto{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Description: user.Description,
		Pic:         user.Pic,
	}
	return userDto, nil
}
