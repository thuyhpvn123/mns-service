package repository

import (
	"github.com/meta-node-blockchain/meta-node-mns/internal/model"
	"gorm.io/gorm"
	"time"
	"github.com/meta-node-blockchain/meta-node-mns/internal/errors"

)

type NameRepository interface {
	Save(name *model.Name) error
	Update(name *model.Name) error
	GetNamesByOwnerAdd(owner string) ([]string,error)
	GetNameFromOwnerAndTokenId(owner string,tokenId uint) (*model.Name,error)
	CheckExpire()([]*model.Name,error)
	GetOwnerByName(fullname string) (string,error)
	GetNameFromOwnerAndLabel(owner string,label string) (*model.Name,error)
}

type nameRepository struct {
	db *gorm.DB
}

func NewNameRepository(db *gorm.DB) NameRepository {
	return &nameRepository{db}
}

func (repo *nameRepository) Save(name *model.Name) error {
	if err := repo.db.Save(name).Error; err != nil {
		return err
	}
	return nil
}

func (repo *nameRepository) Update(name *model.Name) error {
	if err := repo.db.Update("names", name).Error; err != nil {
		return err
	}
	return nil
}
func (repo *nameRepository) GetNamesByOwnerAdd(owner string) ([]string,error) {
	var histories []*model.Name
	var historiesAdd [] string
	currentTime := uint(time.Now().Unix())
	result := repo.db.Model(&model.Name{}).
	Where("owner Like ?", "%"+owner+"%").
	Where("expire_time > ?",currentTime ).
	Find(&histories)
	if result.Error != nil {
        return historiesAdd, result.Error
    }
    for _,v := range histories {
        historiesAdd = append(historiesAdd,v.FullName)
    }
    return historiesAdd, nil
}
func (repo *nameRepository) GetNameFromOwnerAndTokenId(owner string,tokenId uint) (*model.Name,error) {
	var history *model.Name
	result := repo.db.
	Where("owner LIKE ?", "%"+owner+"%").
	Where("token_id = ?",tokenId).
	Find(&history)
	if result.Error != nil {
        return history, result.Error
    }
    return history, nil
}
func (repo *nameRepository) CheckExpire () ([]*model.Name,error){
	var histories []*model.Name
	currentTime := uint(time.Now().Unix())
	result := repo.db.Where("expireTime > ?",currentTime ).Find(&histories)
	if result.Error != nil {
        return histories, result.Error
    }
    return histories, nil
}
func (repo *nameRepository) GetOwnerByName(fullname string) (string,error) {
	var history *model.Name
	currentTime := uint(time.Now().Unix())
	result := repo.db.Model(&model.Name{}).
	Where("full_name Like ?", fullname).
	Where("expire_time > ?",currentTime ).
	Find(&history)
	if result.Error != nil {
        return "", result.Error
    }
	if result.RowsAffected == 0 {
		// No records found, return 404-like response
		return "", errors.ErrNotFound  // You can replace this with a custom error handler if needed
	}
	return history.Owner, nil
}
func (repo *nameRepository) GetNameFromOwnerAndLabel(owner string,label string) (*model.Name,error) {
	var history *model.Name
	fullName := label + ".mtd"
	result := repo.db.
	Where("owner LIKE ?", "%"+owner+"%").
	Where("full_name = ?",fullName).
	Find(&history)
	if result.Error != nil {
        return history, result.Error
    }
    return history, nil
}
