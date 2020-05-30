package group

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/memberrole"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.Provide(new(GroupController))
}

type GroupController struct {
	GroupService *service.GroupService
	UserService  *service.UserService
}

func (g *GroupController) Add() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Name        string `form:"name" binding:"required,min=1"`
			Description string `form:"description" binding:"max=64"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查群组名是否已存在
		if g.GroupService.GetByName(form.Name) != nil {
			return GroupExists
		}

		currentUser := g.UserService.GetCurrentUser(ctx)
		_, err = g.GroupService.Add(currentUser.Id, form.Name, form.Description)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

func (g *GroupController) Update() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Name        string `form:"name" binding:"required,min=1"`
			Description string `form:"description" binding:"max=64"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查群组是否存在
		groupName := ctx.Param("group_name")
		group := g.GroupService.GetByName(groupName)
		if group == nil {
			return GroupNotExists
		}

		// 检查群组名是否被占用
		testGroup := g.GroupService.GetByName(form.Name)
		if testGroup != nil || form.Name != groupName {
			return GroupExists
		}

		// 检查权限
		currentUser := g.UserService.GetCurrentUser(ctx)
		memberInfo := g.GroupService.GetMemberInfo(currentUser.Id, group.Id)
		if memberInfo == nil || memberInfo.Role != memberrole.Owner && memberInfo.Role != memberrole.Admin {
			return PermissionDenied
		}

		err = g.GroupService.Update(group.Id, form.Name, form.Description)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

func (g *GroupController) Delete() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		// 检查群组是否存在
		groupName := ctx.Param("group_name")
		group := g.GroupService.GetByName(groupName)
		if group == nil {
			return GroupNotExists
		}

		// 检查权限
		// 群组所有者才能解散群组
		currentUser := g.UserService.GetCurrentUser(ctx)
		memberInfo := g.GroupService.GetMemberInfo(currentUser.Id, group.Id)
		if memberInfo == nil || memberInfo.Role != memberrole.Owner {
			return PermissionDenied
		}

		err = g.GroupService.Delete(group.Id)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

//func (g *GroupController) Get() gin.HandlerFunc {
//	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
//		// 检查群组是否存在
//		groupName := ctx.Param("group_name")
//		group := g.GroupService.GetByName(groupName)
//		if group == nil {
//			return GroupNotExists
//		}
//		return Success.AddField("group", group)
//	})
//	return handler
//}

func (g *GroupController) ListUserGroup() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		searchWord := ctx.DefaultPostForm("search_word", "")
		currentUser := g.UserService.GetCurrentUser(ctx)
		userGroups := g.GroupService.ListUserGroup(currentUser.Id, searchWord)
		return Success.AddField("user_groups", userGroups)
	})
	return handler
}

//func (g *GroupController) GetUserGroup() gin.HandlerFunc {
//	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
//		return Success
//	})
//	return handler
//}
