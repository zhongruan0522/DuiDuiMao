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

// TierService 档位服务
type TierService struct {
	mode string // dev 或 server
}

// NewTierService 创建档位服务
func NewTierService(mode string) *TierService {
	return &TierService{mode: mode}
}

const cdkCSVPath = "Temp/cdk.csv" // CDK CSV文件路径

// countAvailableCDKs 统计指定档位下可用的CDK数量
func (s *TierService) countAvailableCDKs(tierID int) int {
	if s.mode != config.ModeDev {
		// TODO: 数据库模式实现
		return 0
	}

	// 读取CDK CSV文件
	file, err := os.Open(cdkCSVPath)
	if err != nil {
		// 文件不存在，返回0
		return 0
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0
	}

	count := 0
	for i, record := range records {
		if i == 0 || len(record) < 9 {
			continue // 跳过头部或不完整的行
		}

		// CDK字段：id, tier_id, code, status, order_id, redeemed_by, redeemed_at, created_at, updated_at
		recordTierID, _ := strconv.Atoi(record[1])
		status, _ := strconv.Atoi(record[3])

		// 统计该档位下状态为0（未兑换）的CDK
		if recordTierID == tierID && status == 0 {
			count++
		}
	}

	return count
}

// GetAllTiers 获取所有档位
func (s *TierService) GetAllTiers() ([]model.Tier, error) {
	if s.mode == config.ModeDev {
		return s.readTiersCSV()
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// GetActiveTiers 获取启用的档位（用户端）
func (s *TierService) GetActiveTiers() ([]model.Tier, error) {
	tiers, err := s.GetAllTiers()
	if err != nil {
		return nil, err
	}

	activeTiers := []model.Tier{}
	for _, tier := range tiers {
		if tier.IsActive {
			activeTiers = append(activeTiers, tier)
		}
	}
	return activeTiers, nil
}

// GetTierByID 根据ID获取档位
func (s *TierService) GetTierByID(id int) (*model.Tier, error) {
	tiers, err := s.GetAllTiers()
	if err != nil {
		return nil, err
	}

	for _, tier := range tiers {
		if tier.ID == id {
			return &tier, nil
		}
	}
	return nil, fmt.Errorf("档位不存在")
}

// CreateTier 创建档位（库存自动计算，无需传入）
func (s *TierService) CreateTier(name string, quota, requiredLevel, dailyLimit, sortOrder int, isActive bool) (*model.Tier, error) {
	if s.mode == config.ModeDev {
		return s.createTierCSV(name, quota, requiredLevel, dailyLimit, sortOrder, isActive)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// UpdateTier 更新档位（库存自动计算，无需传入）
func (s *TierService) UpdateTier(id int, name string, quota, requiredLevel, dailyLimit, sortOrder int, isActive bool) (*model.Tier, error) {
	if s.mode == config.ModeDev {
		return s.updateTierCSV(id, name, quota, requiredLevel, dailyLimit, sortOrder, isActive)
	}
	// TODO: 实现数据库版本
	return nil, fmt.Errorf("数据库模式暂未实现")
}

// DeleteTier 删除档位
func (s *TierService) DeleteTier(id int) error {
	if s.mode == config.ModeDev {
		return s.deleteTierCSV(id)
	}
	// TODO: 实现数据库版本
	return fmt.Errorf("数据库模式暂未实现")
}

// UpdateStock 更新库存
func (s *TierService) UpdateStock(id int, delta int) error {
	if s.mode == config.ModeDev {
		return s.updateStockCSV(id, delta)
	}
	// TODO: 实现数据库版本
	return fmt.Errorf("数据库模式暂未实现")
}

// ========== CSV模式实现 ==========

const tierCSVPath = "Temp/tier.csv"

// ensureTierCSV 确保CSV文件存在
func (s *TierService) ensureTierCSV() error {
	dir := filepath.Dir(tierCSVPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	_, statErr := os.Stat(tierCSVPath)
	if os.IsNotExist(statErr) {
		file, createErr := os.Create(tierCSVPath)
		if createErr != nil {
			return createErr
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		// 写入CSV头部
		header := []string{"id", "name", "quota", "required_level", "daily_limit", "stock", "is_active", "sort_order", "created_at", "updated_at"}
		if writeErr := writer.Write(header); writeErr != nil {
			return writeErr
		}
		writer.Flush()
	}
	return nil
}

// readTiersCSV 读取所有档位
func (s *TierService) readTiersCSV() ([]model.Tier, error) {
	if err := s.ensureTierCSV(); err != nil {
		return nil, err
	}

	file, err := os.Open(tierCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tiers := []model.Tier{}
	for i, record := range records {
		if i == 0 || len(record) < 10 {
			continue // 跳过头部或不完整的行
		}

		id, _ := strconv.Atoi(record[0])
		quota, _ := strconv.Atoi(record[2])
		requiredLevel, _ := strconv.Atoi(record[3])
		dailyLimit, _ := strconv.Atoi(record[4])
		// stock字段（record[5]）从CSV读取但会被覆盖
		isActive := record[6] == "true"
		sortOrder, _ := strconv.Atoi(record[7])
		createdAt, _ := time.Parse(time.RFC3339, record[8])
		updatedAt, _ := time.Parse(time.RFC3339, record[9])

		// 自动计算库存：统计该档位下可用的CDK数量
		stock := s.countAvailableCDKs(id)

		tiers = append(tiers, model.Tier{
			ID:            id,
			Name:          record[1],
			Quota:         quota,
			RequiredLevel: requiredLevel,
			DailyLimit:    dailyLimit,
			Stock:         stock, // 使用自动计算的库存
			IsActive:      isActive,
			SortOrder:     sortOrder,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}
	return tiers, nil
}

// writeTiersCSV 写入所有档位
func (s *TierService) writeTiersCSV(tiers []model.Tier) error {
	file, err := os.Create(tierCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入头部
	header := []string{"id", "name", "quota", "required_level", "daily_limit", "stock", "is_active", "sort_order", "created_at", "updated_at"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据
	for _, tier := range tiers {
		record := []string{
			strconv.Itoa(tier.ID),
			tier.Name,
			strconv.Itoa(tier.Quota),
			strconv.Itoa(tier.RequiredLevel),
			strconv.Itoa(tier.DailyLimit),
			strconv.Itoa(tier.Stock),
			strconv.FormatBool(tier.IsActive),
			strconv.Itoa(tier.SortOrder),
			tier.CreatedAt.Format(time.RFC3339),
			tier.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// createTierCSV CSV模式创建档位
func (s *TierService) createTierCSV(name string, quota, requiredLevel, dailyLimit, sortOrder int, isActive bool) (*model.Tier, error) {
	tiers, err := s.readTiersCSV()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newID := 1
	if len(tiers) > 0 {
		newID = tiers[len(tiers)-1].ID + 1
	}

	newTier := model.Tier{
		ID:            newID,
		Name:          name,
		Quota:         quota,
		RequiredLevel: requiredLevel,
		DailyLimit:    dailyLimit,
		Stock:         0, // 新创建的档位库存为0（CSV中写入0，读取时会自动计算）
		IsActive:      isActive,
		SortOrder:     sortOrder,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tiers = append(tiers, newTier)
	if err := s.writeTiersCSV(tiers); err != nil {
		return nil, err
	}

	return &newTier, nil
}

// updateTierCSV CSV模式更新档位
func (s *TierService) updateTierCSV(id int, name string, quota, requiredLevel, dailyLimit, sortOrder int, isActive bool) (*model.Tier, error) {
	tiers, err := s.readTiersCSV()
	if err != nil {
		return nil, err
	}

	var updatedTier *model.Tier
	found := false

	for i := range tiers {
		if tiers[i].ID == id {
			tiers[i].Name = name
			tiers[i].Quota = quota
			tiers[i].RequiredLevel = requiredLevel
			tiers[i].DailyLimit = dailyLimit
			// Stock不更新，保持为0（写入CSV），读取时会自动计算
			tiers[i].IsActive = isActive
			tiers[i].SortOrder = sortOrder
			tiers[i].UpdatedAt = time.Now()
			updatedTier = &tiers[i]
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("档位不存在")
	}

	if err := s.writeTiersCSV(tiers); err != nil {
		return nil, err
	}

	return updatedTier, nil
}

// deleteTierCSV CSV模式删除档位
func (s *TierService) deleteTierCSV(id int) error {
	tiers, err := s.readTiersCSV()
	if err != nil {
		return err
	}

	newTiers := []model.Tier{}
	found := false
	for _, tier := range tiers {
		if tier.ID != id {
			newTiers = append(newTiers, tier)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("档位不存在")
	}

	return s.writeTiersCSV(newTiers)
}

// updateStockCSV CSV模式更新库存（delta为增量，可以是负数）
func (s *TierService) updateStockCSV(id int, delta int) error {
	tiers, err := s.readTiersCSV()
	if err != nil {
		return err
	}

	found := false
	for i := range tiers {
		if tiers[i].ID == id {
			tiers[i].Stock += delta
			if tiers[i].Stock < 0 {
				return fmt.Errorf("库存不足")
			}
			tiers[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("档位不存在")
	}

	return s.writeTiersCSV(tiers)
}
