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

// RedeemLogService 兑换记录服务
type RedeemLogService struct {
	mode string // dev 或 server
}

// NewRedeemLogService 创建兑换记录服务
func NewRedeemLogService(mode string) *RedeemLogService {
	return &RedeemLogService{mode: mode}
}

const redeemLogCSVPath = "Temp/redeem_log.csv"

// CreateRedeemLog 创建兑换记录
func (s *RedeemLogService) CreateRedeemLog(userID, cdkID, tierID int) error {
	if s.mode == config.ModeDev {
		return s.createRedeemLogCSV(userID, cdkID, tierID)
	}
	// TODO: 实现数据库版本
	return fmt.Errorf("数据库模式暂未实现")
}

// GetUserRedeemLogs 获取用户的兑换历史
func (s *RedeemLogService) GetUserRedeemLogs(userID int) ([]model.RedeemLog, error) {
	if s.mode == config.ModeDev {
		return s.getUserRedeemLogsCSV(userID)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// GetAllRedeemLogs 获取所有兑换记录（管理员使用）
func (s *RedeemLogService) GetAllRedeemLogs() ([]model.RedeemLog, error) {
	if s.mode == config.ModeDev {
		return s.getAllRedeemLogsCSV()
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// ========== CSV模式实现 ==========

// ensureRedeemLogCSV 确保兑换记录CSV文件存在
func (s *RedeemLogService) ensureRedeemLogCSV() error {
	dir := filepath.Dir(redeemLogCSVPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	_, statErr := os.Stat(redeemLogCSVPath)
	if os.IsNotExist(statErr) {
		file, createErr := os.Create(redeemLogCSVPath)
		if createErr != nil {
			return createErr
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		// 写入CSV头部
		header := []string{"id", "user_id", "cdk_id", "tier_id", "created_at"}
		if writeErr := writer.Write(header); writeErr != nil {
			return writeErr
		}
		writer.Flush()
	}
	return nil
}

// readRedeemLogsCSV 读取所有兑换记录
func (s *RedeemLogService) readRedeemLogsCSV() ([]model.RedeemLog, error) {
	if err := s.ensureRedeemLogCSV(); err != nil {
		return nil, err
	}

	file, err := os.Open(redeemLogCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	logs := []model.RedeemLog{}
	for i, record := range records {
		if i == 0 || len(record) < 5 {
			continue // 跳过头部或不完整的行
		}

		id, _ := strconv.Atoi(record[0])
		userID, _ := strconv.Atoi(record[1])
		cdkID, _ := strconv.Atoi(record[2])
		tierID, _ := strconv.Atoi(record[3])
		createdAt, _ := time.Parse(time.RFC3339, record[4])

		logs = append(logs, model.RedeemLog{
			ID:        id,
			UserID:    userID,
			CDKID:     cdkID,
			TierID:    tierID,
			CreatedAt: createdAt,
		})
	}
	return logs, nil
}

// writeRedeemLogsCSV 写入所有兑换记录
func (s *RedeemLogService) writeRedeemLogsCSV(logs []model.RedeemLog) error {
	file, err := os.Create(redeemLogCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入头部
	header := []string{"id", "user_id", "cdk_id", "tier_id", "created_at"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据
	for _, log := range logs {
		record := []string{
			strconv.Itoa(log.ID),
			strconv.Itoa(log.UserID),
			strconv.Itoa(log.CDKID),
			strconv.Itoa(log.TierID),
			log.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// createRedeemLogCSV 创建兑换记录（CSV模式）
func (s *RedeemLogService) createRedeemLogCSV(userID, cdkID, tierID int) error {
	logs, err := s.readRedeemLogsCSV()
	if err != nil {
		return err
	}

	// 获取新ID
	newID := 1
	if len(logs) > 0 {
		newID = logs[len(logs)-1].ID + 1
	}

	newLog := model.RedeemLog{
		ID:        newID,
		UserID:    userID,
		CDKID:     cdkID,
		TierID:    tierID,
		CreatedAt: time.Now(),
	}

	logs = append(logs, newLog)
	return s.writeRedeemLogsCSV(logs)
}

// getUserRedeemLogsCSV 获取用户的兑换历史（CSV模式）
func (s *RedeemLogService) getUserRedeemLogsCSV(userID int) ([]model.RedeemLog, error) {
	logs, err := s.readRedeemLogsCSV()
	if err != nil {
		return nil, err
	}

	userLogs := []model.RedeemLog{}
	for _, log := range logs {
		if log.UserID == userID {
			userLogs = append(userLogs, log)
		}
	}

	return userLogs, nil
}

// getAllRedeemLogsCSV 获取所有兑换记录（CSV模式）
func (s *RedeemLogService) getAllRedeemLogsCSV() ([]model.RedeemLog, error) {
	return s.readRedeemLogsCSV()
}
