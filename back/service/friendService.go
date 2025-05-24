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

// GetFriendInfo 根据好友id获取好友信息
func (s *FriendService) GetFriendInfo(friendId int) (dto.FriendInfoDto, error) {
	user, err := s.authRepository.SelectUserById(friendId)
	if user == (models.User{}) {
		return dto.FriendInfoDto{}, err
	}
	friendInfoDto := dto.FriendInfoDto{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Description: user.Description,
		Pic:         user.Pic,
		IsFriend:    1,
	}
	return friendInfoDto, nil
}

// GetFriendList 根据用户id获取好友列表
func (s *FriendService) GetFriendList(userId int) ([]dto.FriendInfoDto, error) {
	friendList, err := s.friendRepository.GetFriendList(userId)
	if err != nil {
		return nil, err
	}
	var friends []dto.FriendInfoDto
	for _, friend := range friendList {
		friendInfoDto := dto.FriendInfoDto{
			Id:          friend.Id,
			Name:        friend.Name,
			Email:       friend.Email,
			Description: friend.Description,
			Pic:         friend.Pic,
			IsFriend:    1,
		}
		friends = append(friends, friendInfoDto)
	}

	return friends, nil
}

// SearchUserByEmail 根据邮箱搜索用户，用于添加好友
func (s *FriendService) SearchUserByEmail(userEmail, searchEmail string) (dto.FriendInfoDto, error) {
	searchUser, err := s.authRepository.SelectUserByEmail(searchEmail)
	if searchUser == (models.User{}) || err != nil {
		return dto.FriendInfoDto{}, err
	}

	// 根据用户id和好友id查询两者目前关系
	// 根据用户邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(userEmail)
	// 查询关系
	relation, err := s.friendRepository.GetFriendRelation(user.Id, searchUser.Id)
	if err != nil {
		return dto.FriendInfoDto{}, err
	}

	var isFriend int
	// 关系不存在
	if relation == (models.Friend{}) {
		isFriend = 2
	} else {
		isFriend = relation.Status
	}

	// 封装返回数据
	friendInfoDto := dto.FriendInfoDto{
		Id:          searchUser.Id,
		Name:        searchUser.Name,
		Email:       searchUser.Email,
		Description: searchUser.Description,
		Pic:         searchUser.Pic,
		IsFriend:    isFriend,
	}
	return friendInfoDto, nil
}
