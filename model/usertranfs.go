package model

import (
	"time"
)

//type UserRecord struct {
//	Usertranfs
//	record
//}
//
//type Usertranfs struct {
//	gorm.Model
//	User
//	tranfs
//}
//
//type tranfs struct {
//	In      int
//	Pay     int
//	balance int
//}
//
//type record struct {
//	zhanji string
//	detail string
//	bei    string
//}

type Transfer struct {
	ID         uint `gorm:"primaryKey"`
	SenderID   uint
	ReceiverID uint
	Amount     float64
	Status     string // "pending", "received", "refunded"
	CreatedAt  time.Time
}
