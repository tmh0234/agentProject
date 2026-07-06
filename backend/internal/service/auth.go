package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo          *repository.Repository
	jwtSecret     []byte
	expireSeconds int64
}

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginResult struct {
	Token   string      `json:"token"`
	Profile UserProfile `json:"profile"`
}

type UserProfile struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Nickname string     `json:"nickname"`
	Status   int        `json:"status"`
	Roles    []RoleInfo `json:"roles"`
}

type RoleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func NewAuthService(db *gorm.DB, secret string, expireSeconds int64) *AuthService {
	if expireSeconds == 0 {
		expireSeconds = 86400
	}
	return &AuthService{repo: repository.New(db), jwtSecret: []byte(secret), expireSeconds: expireSeconds}
}

func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(data), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *AuthService) Login(username, password string) (LoginResult, error) {
	var user model.AdminUser
	err := s.repo.DB().Preload("Roles", "status = ?", 1).Where("username = ? AND status = ?", username, 1).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return LoginResult{}, Unauthorized("用户名或密码错误")
	}
	if err != nil {
		return LoginResult{}, Internal("登录失败", err)
	}
	if !CheckPassword(user.PasswordHash, password) {
		return LoginResult{}, Unauthorized("用户名或密码错误")
	}
	now := time.Now()
	_ = s.repo.DB().Model(&user).Update("last_login_at", &now).Error
	token, err := s.GenerateToken(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, Internal("生成令牌失败", err)
	}
	profile := profileFromUser(user)
	return LoginResult{Token: token, Profile: profile}, nil
}

func (s *AuthService) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.expireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenText string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenText, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, Unauthorized("非法令牌")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, Unauthorized("未登录或登录已过期")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, Unauthorized("未登录或登录已过期")
	}
	return claims, nil
}

func (s *AuthService) Profile(userID uint) (UserProfile, error) {
	var user model.AdminUser
	err := s.repo.DB().Preload("Roles", "status = ?", 1).First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserProfile{}, Unauthorized("用户不存在")
	}
	if err != nil {
		return UserProfile{}, Internal("获取用户信息失败", err)
	}
	return profileFromUser(user), nil
}

func (s *AuthService) Menus(userID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := s.repo.DB().
		Distinct("menus.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN admin_user_roles ON admin_user_roles.role_id = role_menus.role_id").
		Joins("JOIN roles ON roles.id = admin_user_roles.role_id AND roles.deleted_at IS NULL").
		Where("admin_user_roles.admin_user_id = ? AND menus.status = ? AND roles.status = ?", userID, 1, 1).
		Order("menus.parent_id ASC, menus.sort ASC, menus.id ASC").
		Find(&menus).Error
	if err != nil {
		return nil, Internal("获取菜单失败", err)
	}
	return BuildMenuTree(menus), nil
}

func profileFromUser(user model.AdminUser) UserProfile {
	roles := make([]RoleInfo, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, RoleInfo{ID: role.ID, Name: role.Name, Code: role.Code})
	}
	return UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Status:   user.Status,
		Roles:    roles,
	}
}

func BuildMenuTree(items []model.Menu) []model.Menu {
	byParent := map[uint][]model.Menu{}
	for _, item := range items {
		item.Children = nil
		byParent[item.ParentID] = append(byParent[item.ParentID], item)
	}
	var attach func(parentID uint) []model.Menu
	attach = func(parentID uint) []model.Menu {
		children := byParent[parentID]
		for i := range children {
			children[i].Children = attach(children[i].ID)
		}
		return children
	}
	return attach(0)
}
