package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zhongruan0522/DuiDuiMao/internal/config"
	"github.com/zhongruan0522/DuiDuiMao/internal/model"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// CDKService CDK服务
type CDKService struct {
	mode string // dev 或 server
}

// NewCDKService 创建CDK服务
func NewCDKService(mode string) *CDKService {
	return &CDKService{mode: mode}
}

const cdkCSVPath = "Temp/cdk.csv"

// ImportCDKsRequest 批量导入CDK请求
type ImportCDKsRequest struct {
	TierID int      `json:"tier_id"` // 所属档位ID
	Codes  []string `json:"codes"`   // CDK列表（一行一个）
}

// ImportCDKsResult 批量导入CDK结果
type ImportCDKsResult struct {
	SuccessCount int      `json:"success_count"` // 成功导入数量
	FailedCount  int      `json:"failed_count"`  // 失败数量
	FailedCodes  []string `json:"failed_codes"`  // 失败的CDK列表（重复的）
}

// BatchImportCDKs 批量导入CDK
func (s *CDKService) BatchImportCDKs(tierID int, codes []string) (*ImportCDKsResult, error) {
	if s.mode == config.ModeDev {
		return s.batchImportCDKsCSV(tierID, codes)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// GetCDKs 获取CDK列表（支持按档位和状态筛选）
func (s *CDKService) GetCDKs(tierID *int, status *int) ([]model.CDK, error) {
	if s.mode == config.ModeDev {
		return s.getCDKsCSV(tierID, status)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// RevokeCDK 作废CDK
func (s *CDKService) RevokeCDK(id int) error {
	if s.mode == config.ModeDev {
		return s.revokeCDKCSV(id)
	}
	// TODO: 实现数据库版本
	return fmt.Errorf("数据库模式暂未实现")
}

// GetAvailableCDKByTierID 获取指定档位的一个可用CDK（用于兑换）
func (s *CDKService) GetAvailableCDKByTierID(tierID int) (*model.CDK, error) {
	if s.mode == config.ModeDev {
		return s.getAvailableCDKByTierIDCSV(tierID)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// MarkCDKAsRedeemed 标记CDK为已兑换（兑换时使用）
func (s *CDKService) MarkCDKAsRedeemed(cdkID, userID int) error {
	if s.mode == config.ModeDev {
		return s.markCDKAsRedeemedCSV(cdkID, userID)
	}
	// TODO: 实现数据库版本
	return fmt.Errorf("数据库模式暂未实现")
}

// ========== CSV模式实现 ==========

// ensureCDKCSV 确保CDK CSV文件存在
func (s *CDKService) ensureCDKCSV() error {
	dir := filepath.Dir(cdkCSVPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	_, statErr := os.Stat(cdkCSVPath)
	if os.IsNotExist(statErr) {
		file, createErr := os.Create(cdkCSVPath)
		if createErr != nil {
			return createErr
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		// 写入CSV头部
		// 字段：id, tier_id, code, status, order_id, redeemed_by, redeemed_at, created_at, updated_at
		header := []string{"id", "tier_id", "code", "status", "order_id", "redeemed_by", "redeemed_at", "created_at", "updated_at"}
		if writeErr := writer.Write(header); writeErr != nil {
			return writeErr
		}
		writer.Flush()
	}
	return nil
}

// readCDKsCSV 读取所有CDK
func (s *CDKService) readCDKsCSV() ([]model.CDK, error) {
	if err := s.ensureCDKCSV(); err != nil {
		return nil, err
	}

	file, err := os.Open(cdkCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	cdks := []model.CDK{}
	for i, record := range records {
		if i == 0 || len(record) < 9 {
			continue // 跳过头部或不完整的行
		}

		id, _ := strconv.Atoi(record[0])
		tierID, _ := strconv.Atoi(record[1])
		status, _ := strconv.Atoi(record[3])
		orderID, _ := strconv.Atoi(record[4])
		redeemedBy, _ := strconv.Atoi(record[5])
		redeemedAt, _ := time.Parse(time.RFC3339, record[6])
		createdAt, _ := time.Parse(time.RFC3339, record[7])
		updatedAt, _ := time.Parse(time.RFC3339, record[8])

		cdks = append(cdks, model.CDK{
			ID:         id,
			TierID:     tierID,
			Code:       record[2], // 已加密存储
			Status:     status,
			OrderID:    orderID,
			RedeemedBy: redeemedBy,
			RedeemedAt: redeemedAt,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}
	return cdks, nil
}

// writeCDKsCSV 写入所有CDK
func (s *CDKService) writeCDKsCSV(cdks []model.CDK) error {
	file, err := os.Create(cdkCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入头部
	header := []string{"id", "tier_id", "code", "status", "order_id", "redeemed_by", "redeemed_at", "created_at", "updated_at"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据
	for _, cdk := range cdks {
		redeemedAtStr := ""
		if !cdk.RedeemedAt.IsZero() {
			redeemedAtStr = cdk.RedeemedAt.Format(time.RFC3339)
		}

		record := []string{
			strconv.Itoa(cdk.ID),
			strconv.Itoa(cdk.TierID),
			cdk.Code, // 已加密
			strconv.Itoa(cdk.Status),
			strconv.Itoa(cdk.OrderID),
			strconv.Itoa(cdk.RedeemedBy),
			redeemedAtStr,
			cdk.CreatedAt.Format(time.RFC3339),
			cdk.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// batchImportCDKsCSV 批量导入CDK（CSV模式）
func (s *CDKService) batchImportCDKsCSV(tierID int, codes []string) (*ImportCDKsResult, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("CDK列表不能为空")
	}

	cdks, err := s.readCDKsCSV()
	if err != nil {
		return nil, err
	}

	// 获取新ID起始值
	newID := 1
	if len(cdks) > 0 {
		newID = cdks[len(cdks)-1].ID + 1
	}

	// 构建已存在CDK的映射（用于检测重复）
	existingCodes := make(map[string]bool)
	for _, cdk := range cdks {
		// 解密后比对（假设存储时已加密）
		decryptedCode, decErr := util.DoubleDecode(cdk.Code)
		if decErr == nil {
			existingCodes[decryptedCode] = true
		}
	}

	result := &ImportCDKsResult{
		SuccessCount: 0,
		FailedCount:  0,
		FailedCodes:  []string{},
	}

	now := time.Now()

	// 批量添加CDK
	for _, code := range codes {
		trimmedCode := strings.TrimSpace(code)
		if trimmedCode == "" {
			continue // 跳过空行
		}

		// 检查重复
		if existingCodes[trimmedCode] {
			result.FailedCount++
			result.FailedCodes = append(result.FailedCodes, trimmedCode)
			continue
		}

		// 加密存储
		encryptedCode := util.DoubleEncode(trimmedCode)

		newCDK := model.CDK{
			ID:         newID,
			TierID:     tierID,
			Code:       encryptedCode,
			Status:     0, // 0=未兑换
			OrderID:    0,
			RedeemedBy: 0,
			RedeemedAt: time.Time{}, // 零值
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		cdks = append(cdks, newCDK)
		existingCodes[trimmedCode] = true
		result.SuccessCount++
		newID++
	}

	// 写回CSV
	if err := s.writeCDKsCSV(cdks); err != nil {
		return nil, err
	}

	return result, nil
}

// getCDKsCSV 获取CDK列表（CSV模式，支持筛选）
func (s *CDKService) getCDKsCSV(tierID *int, status *int) ([]model.CDK, error) {
	cdks, err := s.readCDKsCSV()
	if err != nil {
		return nil, err
	}

	// 筛选
	filtered := []model.CDK{}
	for _, cdk := range cdks {
		// 按档位筛选
		if tierID != nil && cdk.TierID != *tierID {
			continue
		}
		// 按状态筛选
		if status != nil && cdk.Status != *status {
			continue
		}
		filtered = append(filtered, cdk)
	}

	return filtered, nil
}

// revokeCDKCSV 作废CDK（CSV模式）
func (s *CDKService) revokeCDKCSV(id int) error {
	cdks, err := s.readCDKsCSV()
	if err != nil {
		return err
	}

	found := false
	for i := range cdks {
		if cdks[i].ID == id {
			if cdks[i].Status == 2 {
				return fmt.Errorf("CDK已被兑换，无法作废")
			}
			cdks[i].Status = 3 // 3=已作废
			cdks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("CDK不存在")
	}

	return s.writeCDKsCSV(cdks)
}

// getAvailableCDKByTierIDCSV 获取指定档位的一个可用CDK（CSV模式）
func (s *CDKService) getAvailableCDKByTierIDCSV(tierID int) (*model.CDK, error) {
	cdks, err := s.readCDKsCSV()
	if err != nil {
		return nil, err
	}

	for i := range cdks {
		if cdks[i].TierID == tierID && cdks[i].Status == 0 {
			return &cdks[i], nil
		}
	}

	return nil, fmt.Errorf("该档位暂无可用CDK")
}

// markCDKAsRedeemedCSV 标记CDK为已兑换（CSV模式）
func (s *CDKService) markCDKAsRedeemedCSV(cdkID, userID int) error {
	cdks, err := s.readCDKsCSV()
	if err != nil {
		return err
	}

	found := false
	for i := range cdks {
		if cdks[i].ID == cdkID {
			if cdks[i].Status != 0 {
				return fmt.Errorf("CDK已被兑换或作废")
			}
			cdks[i].Status = 2 // 2=已兑换
			cdks[i].RedeemedBy = userID
			cdks[i].RedeemedAt = time.Now()
			cdks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("CDK不存在")
	}

	return s.writeCDKsCSV(cdks)
}
