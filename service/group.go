package service

import (
	"fmt"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/memberrole"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	inject.Provide(new(GroupService))
}

type GroupService struct {
	UserService *UserService
}

func (g *GroupService) Add(ownerUserId int, name, description string) (*model.Group, error) {
	var err error
	now := time.Now()

	// 开事务
	tx := dao.DB.Begin()
	// 添加群组
	group := model.Group{
		OwnerUserId: ownerUserId,
		Name:        name,
		Description: description,
		UpdateTime:  now,
		CreateTime:  now,
	}
	if err = tx.Create(&group).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// 添加群组成员信息
	groupMember := model.GroupMember{
		GroupId:    group.Id,
		UserId:     ownerUserId,
		Role:       memberrole.Owner,
		UpdateTime: now,
		CreateTime: now,
	}
	if err = tx.Create(&groupMember).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *GroupService) Update(groupId int, name, description string) error {
	now := time.Now()
	group := model.Group{Id: groupId}

	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"update_time": now,
	}
	if err := dao.DB.Model(&group).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (g *GroupService) GetByName(name string) *model.Group {
	var group model.Group
	if err := dao.DB.First(&group, "`name`=?", name).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &group
}

func (g *GroupService) Delete(groupId int) error {
	var err error
	tx := dao.DB.Begin()
	// 删除成员
	if err = tx.Where("`group_id`=?", groupId).Delete(model.GroupMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO: 删除群组分享

	// 删除群组
	if err = tx.Where("`id`=?", groupId).Delete(model.Group{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (g *GroupService) ListUserGroup(userId int, searchWord string) []*model.UserGroup {
	members := g.ListMembersByUserId(userId)

	// 查询需要的群组
	groupIds := make([]int, 0, len(members))
	for _, member := range members {
		groupIds = append(groupIds, member.GroupId)
	}
	var groups []*model.Group
	if searchWord == "" {
		if err := dao.DB.Where("`id` IN (?)", groupIds).Find(&groups).Error; err != nil {
			panic(err)
		}
	} else {
		matchString := fmt.Sprintf("%s%%", searchWord)
		if err := dao.DB.Where("`id` IN (?) AND `name` LIKE ?", groupIds, matchString).Find(&groups).Error; err != nil {
			panic(err)
		}
	}
	// 构造群组表，方便下面映射
	groupMap := make(map[int]*model.Group)
	for _, group := range groupMap {
		groupMap[group.Id] = group
	}

	// 构造UserGroup列表
	user := g.UserService.GetById(userId)
	userGroups := make([]*model.UserGroup, 0)
	for _, member := range members {
		group := groupMap[member.GroupId]
		if group == nil {
			continue
		}
		userGroups = append(userGroups, &model.UserGroup{
			Group: group,
			User:  user,
			Role:  member.Role,
		})
	}
	return userGroups
}

func (g *GroupService) GetMemberInfo(userId, groupId int) *model.GroupMember {
	var groupMember model.GroupMember
	if err := dao.DB.First(&groupMember, "`user_id`=? AND `group_id`=?", userId, groupId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &groupMember
}

func (g *GroupService) ListMembersByUserId(userId int) []*model.GroupMember {
	var members []*model.GroupMember
	if err := dao.DB.Where("`user_id`=?", userId).Find(&members).Error; err != nil {
		panic(err)
	}
	return members
}

func (g *GroupService) ListMembersByGroupId(groupId int) []*model.GroupMember {
	var members []*model.GroupMember
	if err := dao.DB.Where("`group_id`=?", groupId).Find(&members).Error; err != nil {
		panic(err)
	}
	return members
}
