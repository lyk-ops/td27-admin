package sysTool

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"server/global"
	commonReq "server/model/common/request"
	modelSysTool "server/model/sysTool"
	sysToolReq "server/model/sysTool/request"
	"server/utils"
)

type CronService struct{}

// GetCronList 分页获取cron
func (cs *CronService) GetCronList(pageInfo commonReq.PageInfo) (cronModelList []modelSysTool.CronModel, total int64, err error) {
	db := global.TD27_DB.Model(&modelSysTool.CronModel{})

	// 计算记录数量
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count err %v", err)
	}

	// 分页
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	if pageInfo.PageSize > 0 && pageInfo.Page > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&cronModelList).Error
	return
}

// AddCron 添加cron
func (cs *CronService) AddCron(cronModel *modelSysTool.CronModel) (*modelSysTool.CronModel, error) {
	// 开启cron
	if cronModel.Open {
		entryId, err := utils.AddJob(cronModel)
		if err != nil {
			return nil, err
		} else {
			cronModel.EntryId = entryId
		}
	}
	err := global.TD27_DB.Create(cronModel).Error
	return cronModel, err
}

// DeleteCron 删除cron
func (cs *CronService) DeleteCron(id uint) error {
	var cronModel modelSysTool.CronModel
	if errors.Is(global.TD27_DB.Where("id = ?", id).First(&cronModel).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录未找到")
	}
	// 删除定时任务
	global.TD27_CRON.Remove(cron.EntryID(cronModel.EntryId))
	// 删除数据库记录
	return global.TD27_DB.Unscoped().Delete(&cronModel).Error
}

// DeleteCronByIds 批量删除cron
func (cs *CronService) DeleteCronByIds(ids []uint) error {
	var cronModels []modelSysTool.CronModel
	global.TD27_DB.Find(&cronModels, ids)
	// 删除定时任务
	for _, value := range cronModels {
		global.TD27_CRON.Remove(cron.EntryID(value.EntryId))
	}
	// 删除数据库记录
	return global.TD27_DB.Unscoped().Delete(&cronModels).Error
}

// EditCron 编辑cron
func (cs *CronService) EditCron(cronReq *sysToolReq.CronReq) (*modelSysTool.CronModel, error) {
	var cronModel modelSysTool.CronModel
	if errors.Is(global.TD27_DB.Where("id = ?", cronReq.ID).First(&cronModel).Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("记录未找到")
	}
	// 拼接
	cronModel.Name = cronReq.Name
	cronModel.Method = cronReq.Method
	cronModel.Expression = cronReq.Expression
	cronModel.Strategy = cronReq.Strategy
	// params 拼接
	for _, v := range cronModel.ExtraParams.TableInfo {
		cronModel.ExtraParams.TableInfo = append(cronModel.ExtraParams.TableInfo, v)
	}
	cronModel.ExtraParams.Command = cronReq.ExtraParams.Command
	cronModel.Comment = cronReq.Comment
	if cronReq.Open {
		//utils.IsContain(utils.GetEntries(), cronModel.EntryId)
		if cronModel.EntryId <= 0 {
			entryId, err := utils.AddJob(&cronModel)
			if err != nil {
				global.TD27_LOG.Error("Add cron", zap.Error(err))
			}
			cronModel.EntryId = entryId
		}
	} else {
		if cronModel.EntryId > 0 {
			global.TD27_CRON.Remove(cron.EntryID(cronModel.EntryId))
			cronModel.EntryId = 0
		}
	}
	err := global.TD27_DB.Save(&cronModel).Error
	return &cronModel, err
}

// SwitchOpen 切换cron活跃状态
func (cs *CronService) SwitchOpen(id uint, open bool) (err error) {
	var cronModel modelSysTool.CronModel
	if errors.Is(global.TD27_DB.Where("id = ?", id).First(&cronModel).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录未找到")
	}

	if open {
		if cronModel.EntryId <= 0 {
			entryId, err := utils.AddJob(&cronModel)
			if err != nil {
				global.TD27_LOG.Error("Add cron", zap.Error(err))
				return err
			}
			return global.TD27_DB.Model(&cronModel).Updates(map[string]interface{}{"open": true, "entryId": entryId}).Error
		}
	} else {
		if cronModel.EntryId != 0 {
			global.TD27_CRON.Remove(cron.EntryID(cronModel.EntryId))
			return global.TD27_DB.Model(&cronModel).Updates(map[string]interface{}{"open": false, "entryId": 0}).Error
		}
	}

	return
}
