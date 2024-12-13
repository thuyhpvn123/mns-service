package usecase

import (
	"github.com/meta-node-blockchain/meta-node-mns/internal/model"
	"github.com/meta-node-blockchain/meta-node-mns/internal/repository"
)
type NameUsecase interface {
	Save(name *model.Name) error
	Update(name *model.Name) error
	GetNamesByOwnerAdd(owner string) ([]string,error)
	GetNameFromOwnerAndTokenId(owner string,tokenId uint) (*model.Name,error)
	CheckExpire() ([]*model.Name,error)
	GetOwnerByName(name string) (string,error)
	GetNameFromOwnerAndLabel(owner string,label string) (*model.Name,error)
}

type nameUsecase struct {
	nameRepo repository.NameRepository
}
func NewNameUsecase(nameRepo repository.NameRepository) NameUsecase {
	return &nameUsecase{nameRepo}
}

func (usecase *nameUsecase) Save(name *model.Name) error {
	return usecase.nameRepo.Save(name)
}

func (usecase *nameUsecase) Update(name *model.Name) error {
	return usecase.nameRepo.Update(name)
}
func (usecase *nameUsecase) GetNamesByOwnerAdd(owner string) ([]string,error) {
	return usecase.nameRepo.GetNamesByOwnerAdd(owner)
}

func (usecase *nameUsecase) GetNameFromOwnerAndTokenId(owner string,tokenId uint) (*model.Name,error) {
	return usecase.nameRepo.GetNameFromOwnerAndTokenId(owner,tokenId)
}
func (usecase *nameUsecase) CheckExpire() ([]*model.Name,error) {
	return usecase.nameRepo.CheckExpire()
}
func (usecase *nameUsecase) GetOwnerByName(name string) (string,error) {
	return usecase.nameRepo.GetOwnerByName(name)
}
func (usecase *nameUsecase) GetNameFromOwnerAndLabel(owner string,label string) (*model.Name,error) {
	return usecase.nameRepo.GetNameFromOwnerAndLabel(owner,label)
}
