package model

type Name struct {
	ID 			uint64 `gorm:"primaryKey" json:"id"`
	TokenID 	uint `json:"tokenId" gorm:"column:token_id"`   //uint256(hash(labelhash,hash(mtd))
	ExpireTime  uint `json:"expireTime" gorm:"column:expire_time"`
	Owner 		string `json:"owner" gorm:"column:owner"`
	FullName 	string `json:"fullname" gorm:"column:full_name"`
}