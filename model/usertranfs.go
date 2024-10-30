package model

import "gorm.io/gorm"

type Usertranfs struct {
	gorm.Model
	User
	tranfs
}

type tranfs struct {
	In      int
	Pay     int
	balance int
}

type record struct {
	zhanji string
	detail string
	bei    string
}

type UserRecord struct {
	Usertranfs
	record
}
