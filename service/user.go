package service

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/usertype"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"time"
)

const (
	DefaultSystemAdminUsername = "admin"
	DefaultUserPassword        = "123456"
	TokenCookieKey             = "TOKEN"
	JwtSignatureKey            = "v&g094gydc9_svpk@+02@jl5loa-)sk37g@lffbur-pj#5#=3!"
	CurrentUserContextName     = "currentUser"
)

var (
	jwtSignatureKey     = []byte(JwtSignatureKey)
	jwtSignatureKeyFunc = func(token *jwt.Token) (interface{}, error) {
		return jwtSignatureKey, nil
	}
)

func init() {
	inject.Provide(new(UserService))
}

type UserService struct{}

// GenerateToken 生成token
func (u *UserService) GenerateToken(user *model.User) (string, error) {
	userJsonBytes, err := jsoniter.Marshal(user)
	if err != nil {
		return "", err
	}
	userInfoJson := string(userJsonBytes)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["userInfo"] = userInfoJson
	token.Claims = claims

	tokenString, err := token.SignedString(jwtSignatureKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CalculateHashPassword 计算哈希密码
func (u *UserService) CalculateHashPassword(password string) string {
	sha1Result := sha1.Sum([]byte(password))
	sha256Result := sha256.Sum256(sha1Result[0:])
	return hex.EncodeToString(sha256Result[0:])
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(ctx *gin.Context) *model.User {
	data, exists := ctx.Get(CurrentUserContextName)
	if !exists {
		return nil
	}
	user, ok := data.(*model.User)
	if !ok {
		return nil
	}
	return user
}
func (u *UserService) GetCurrentUser(ctx *gin.Context) *model.User {
	return GetCurrentUser(ctx)
}

// GetByUsername 根据用户名获取用户信息
func (u *UserService) GetByUsername(username string) *model.User {
	var user model.User
	if err := dao.DB.First(&user, "username=?", username).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &user
}

// GetById 根据用户ID获取用户信息
func (u *UserService) GetById(userId int) *model.User {
	var user model.User
	if err := dao.DB.First(&user, userId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &user
}

//// ListAll 获取所有用户信息
//func (u *UserService) ListAll() []*model.User {
//	var selectResult []entity.User
//	if err := dao.DB.Find(&selectResult).Error; err != nil {
//		panic(err)
//	}
//
//	users := make([]*model.User, len(selectResult))
//	for i, user := range selectResult {
//		users[i] = new(model.User)
//		users[i].FromEntity(&user)
//	}
//	return users
//}
//
//// ListByIds 获取指定用户ID列表的用户信息
//func (u *UserService) ListByIds(userIds []int) []*model.User {
//	panic("implement me")
//}
//
//// ListNormalUsers 获取所有普通用户信息
//func (u *UserService) ListNormalUsers() []*model.User {
//	var selectResult []entity.User
//	if err := dao.DB.Where("type=?", 2).Find(&selectResult).Error; err != nil {
//		panic(err)
//	}
//
//	users := make([]*model.User, len(selectResult))
//	for i, user := range selectResult {
//		users[i] = new(model.User)
//		users[i].FromEntity(&user)
//	}
//	return users
//}
//
//// Search 根据关键词搜索用户名与昵称
//func (u *UserService) Search(searchWord string) []*model.User {
//	var selectResult []entity.User
//	if err := dao.DB.Where("`username`=? or `nickname`=?", searchWord, searchWord).Find(&selectResult).Error; err != nil {
//		panic(err)
//	}
//
//	users := make([]*model.User, len(selectResult))
//	for i, user := range selectResult {
//		users[i] = new(model.User)
//		users[i].FromEntity(&user)
//	}
//	return users
//}
//
//// ExistsAll 判断是否所有指定ID的用户都存在
//func (u *UserService) ExistsAll(userIds []int) bool {
//	panic("implement me")
//}

// AddDefaultSystemAdmin 添加默认系统管理员
func (u *UserService) AddDefaultSystemAdmin() (*model.User, error) {
	now := time.Now()
	user := model.User{
		Type:       usertype.SystemAdmin,
		Username:   DefaultSystemAdminUsername,
		Password:   u.CalculateHashPassword(DefaultUserPassword),
		Nickname:   "",
		Remark:     "",
		UpdateTime: now,
		CreateTime: now,
	}
	if err := dao.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

//// Add 添加新用户
//func (u *UserService) Add(username, remark, homeDirectory string) (*model.User, error) {
//	now := time.Now()
//	userEntity := entity.User{
//		Type:          model.UserTypeNormal.Value(),
//		Username:      username,
//		Password:      u.CalculateHashPassword(DefaultUserPassword),
//		Nickname:      "",
//		Remark:        remark,
//		HomeDirectory: homeDirectory,
//		UpdateTime:    now,
//		CreateTime:    now,
//	}
//	if err := dao.DB.Create(&userEntity).Error; err != nil {
//		return nil, err
//	}
//	return (&model.User{}).FromEntity(&userEntity), nil
//}
//
//// DeleteById 删除用户
//func (u *UserService) DeleteById(userId int) error {
//	// TODO: 删除用户分享的
//
//	// TODO: 删除分享给用户的
//
//	// TODO: 删除拥有的群组
//
//	// TODO: 退出其他群组
//
//	// TODO: 通知更新share_record表的target_content
//
//	// TODO: 删除用户
//	panic("implement me")
//}
//
//// Update 更新指定用户信息
//func (u *UserService) Update(userId int, nickname, remark string) (*model.User, error) {
//	user := entity.User{Id: userId}
//
//	updates := map[string]interface{}{
//		"nickname": nickname,
//		"remark":   remark,
//	}
//	if err := dao.DB.Model(&user).Updates(updates).Error; err != nil {
//		return nil, err
//	}
//	return u.GetById(userId), nil
//}
//
//// UpdateUserPassword 更新指定用户的密码
//func (u *UserService) UpdateUserPassword(userId int, password string) error {
//	user := entity.User{Id: userId}
//
//	updates := map[string]interface{}{
//		"password": u.CalculateHashPassword(password),
//	}
//	if err := dao.DB.Model(&user).Updates(updates).Error; err != nil {
//		return err
//	}
//	return nil
//}

// AddLoginRecord 添加用户登录记录
func (u *UserService) AddLoginRecord(user *model.User, source string) error {
	userLoginRecord := model.UserLoginRecord{
		UserId:     user.Id,
		Source:     source,
		CreateTime: time.Now(),
	}
	if err := dao.DB.Create(&userLoginRecord).Error; err != nil {
		return err
	}
	return nil
}
