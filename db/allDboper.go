package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"main/models"
	"time"
)

type Repository struct {
	Db *gorm.DB
}

func GetConnection(dbSetting string) (*Repository, error) {
	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		return nil, err
	}
	return &Repository{Db: db}, nil
}

func (r *Repository) CreateProfile(profiles *models.CommissionProfiles) (int64, error) {
	if err := r.Db.Create(&profiles).Error; err != nil {
		return 0, nil
	}
	return profiles.Id, nil
}

func (r *Repository) CreateRules(rules []models.CommissionRules, profileId int64) error {
	for _, rule := range rules {
		*rule.Active = true
		rule.ProfileId = profileId
		err := r.Db.Create(&rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdateProfile(profiles *models.CommissionProfiles) (*models.CommissionProfiles, error) {
	err := r.Db.Model(&models.CommissionProfiles{}).Update(&profiles).Scan(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *Repository) DeleteProfile(profiles models.CommissionProfiles) error {
	if err := r.Db.Model(&models.CommissionProfiles{}).Where("id = ?", profiles.Id).UpdateColumns(map[string]interface{}{
		"active":     false,
		"deleted_at": profiles.UpdatedAt,
		"updated_by": profiles.UpdatedBy,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateRules(rule *models.CommissionRules) (*models.CommissionRules, error) {
	if err := r.Db.Model(&models.CommissionRules{}).Where("id = ?", rule.Id).
		Update(&rule).Scan(&rule).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

func (r *Repository) DeleteRule(rules *models.CommissionRules) error {
	if err := r.Db.Model(&models.CommissionRules{}).Where("id = ?", rules.Id).
		UpdateColumns(map[string]interface{}{
			"active":     false,
			"deleted_at": time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateProfileRules(profileId, updaterID int64) error {
	if err := r.Db.Model(&models.CommissionProfiles{}).Where("id = ?", profileId).UpdateColumns(map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": updaterID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateProfileAndRules(request *models.ProfileCreatRequest) error {
	tx := r.Db.Begin()
	err := tx.Model(&models.CommissionProfiles{}).Update(&request.Profile).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, rule := range request.Rules {
		rule.UpdatedAt = new(time.Time)
		*rule.UpdatedAt = time.Now()
		if err := tx.Save(&rule).Error; err != nil {
			tx.Rollback()
			return err
		}
		//if err := tx.Model(&models.CommissionRules{}).Where("id = ?", rule.Id).
		//	UpdateColumns(map[string]interface{}{
		//		"index":       rule.Index,
		//		"start_range": rule.StartRange,
		//		"end_range":   rule.EndRange,
		//		"value":       rule.Value,
		//		"type_id":     rule.TypeId,
		//		"active":      rule.Active,
		//		"updated_at":  time.Now(),
		//	}).Error; err != nil {
		//	tx.Rollback()
		//	return err
		//}
	}
	tx.Commit()
	return nil
}

func (r *Repository) DeleteProfileAndRules(request *models.ProfileCreatRequest) error {
	tx := r.Db.Begin()
	if err := tx.Model(&models.CommissionProfiles{}).Where("id = ?", request.Profile.Id).UpdateColumns(map[string]interface{}{
		"active":     false,
		"deleted_at": request.Profile.DeletedAt,
		"updated_by": request.Profile.UpdatedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.CommissionRules{}).Where("profile_id = ?", request.Profile.Id).
		UpdateColumns(map[string]interface{}{
			"active":     false,
			"deleted_at": time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
