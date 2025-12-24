package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/zhongruan0522/DuiDuiMao/internal/config"
	"github.com/zhongruan0522/DuiDuiMao/internal/model"
)

// UserService 用户服务
type UserService struct {
	mode string // dev 或 server
}

// NewUserService 创建用户服务
func NewUserService(mode string) *UserService {
	return &UserService{mode: mode}
}

// CreateOrUpdateUser 创建或更新用户（OAuth登录后调用）
func (s *UserService) CreateOrUpdateUser(linuxDoID int, username, name string, trustLevel int, isAdmin bool) (*model.User, error) {
	if s.mode == config.ModeDev {
		return s.createOrUpdateUserCSV(linuxDoID, username, name, trustLevel, isAdmin)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// GetUserByLinuxDoID 根据LinuxDo ID获取用户
func (s *UserService) GetUserByLinuxDoID(linuxDoID int) (*model.User, error) {
	if s.mode == config.ModeDev {
		return s.getUserByLinuxDoIDCSV(linuxDoID)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int) (*model.User, error) {
	if s.mode == config.ModeDev {
		return s.getUserByIDCSV(id)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// ========== CSV模式实现 ==========

const userCSVPath = "Temp/user.csv"

// ensureUserCSV 确保CSV文件存在
func (s *UserService) ensureUserCSV() error {
	dir := filepath.Dir(userCSVPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(userCSVPath); os.IsNotExist(err) {
		file, err := os.Create(userCSVPath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		// 写入CSV头部
		header := []string{"id", "linux_do_id", "username", "name", "trust_level", "is_admin", "created_at", "updated_at"}
		if err := writer.Write(header); err != nil {
			return err
		}
		writer.Flush()
	}
	return nil
}

// readUsersCSV 读取所有用户
func (s *UserService) readUsersCSV() ([]model.User, error) {
	if err := s.ensureUserCSV(); err != nil {
		return nil, err
	}

	file, err := os.Open(userCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	users := []model.User{}
	for i, record := range records {
		if i == 0 || len(record) < 8 {
			continue // 跳过头部或不完整的行
		}

		id, _ := strconv.Atoi(record[0])
		linuxDoID, _ := strconv.Atoi(record[1])
		trustLevel, _ := strconv.Atoi(record[4])
		isAdmin := record[5] == "true"
		createdAt, _ := time.Parse(time.RFC3339, record[6])
		updatedAt, _ := time.Parse(time.RFC3339, record[7])

		users = append(users, model.User{
			ID:         id,
			LinuxDoID:  linuxDoID,
			Username:   record[2],
			Name:       record[3],
			TrustLevel: trustLevel,
			IsAdmin:    isAdmin,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}
	return users, nil
}

// writeUsersCSV 写入所有用户
func (s *UserService) writeUsersCSV(users []model.User) error {
	file, err := os.Create(userCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入头部
	header := []string{"id", "linux_do_id", "username", "name", "trust_level", "is_admin", "created_at", "updated_at"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据
	for _, user := range users {
		record := []string{
			strconv.Itoa(user.ID),
			strconv.Itoa(user.LinuxDoID),
			user.Username,
			user.Name,
			strconv.Itoa(user.TrustLevel),
			strconv.FormatBool(user.IsAdmin),
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// createOrUpdateUserCSV CSV模式创建或更新用户
func (s *UserService) createOrUpdateUserCSV(linuxDoID int, username, name string, trustLevel int, isAdmin bool) (*model.User, error) {
	users, err := s.readUsersCSV()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var targetUser *model.User

	// 查找是否已存在
	found := false
	for i := range users {
		if users[i].LinuxDoID == linuxDoID {
			// 更新现有用户
			users[i].Username = username
			users[i].Name = name
			users[i].TrustLevel = trustLevel
			users[i].IsAdmin = isAdmin
			users[i].UpdatedAt = now
			targetUser = &users[i]
			found = true
			break
		}
	}

	if !found {
		// 创建新用户
		newID := 1
		if len(users) > 0 {
			newID = users[len(users)-1].ID + 1
		}

		newUser := model.User{
			ID:         newID,
			LinuxDoID:  linuxDoID,
			Username:   username,
			Name:       name,
			TrustLevel: trustLevel,
			IsAdmin:    isAdmin,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		users = append(users, newUser)
		targetUser = &newUser
	}

	if err := s.writeUsersCSV(users); err != nil {
		return nil, err
	}

	return targetUser, nil
}

// getUserByLinuxDoIDCSV CSV模式根据LinuxDo ID获取用户
func (s *UserService) getUserByLinuxDoIDCSV(linuxDoID int) (*model.User, error) {
	users, err := s.readUsersCSV()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.LinuxDoID == linuxDoID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("用户不存在")
}

// getUserByIDCSV CSV模式根据ID获取用户
func (s *UserService) getUserByIDCSV(id int) (*model.User, error) {
	users, err := s.readUsersCSV()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("用户不存在")
}
