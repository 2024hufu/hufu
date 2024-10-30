package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 按照依赖关系顺序进行迁移
	err = db.AutoMigrate(
		&Wallet{},
		&WalletKey{},
		&Transaction{},
		&EncryptedTransaction{},
		&DesensitizedTransaction{},
		&AbnormalTransaction{},
	)
	if err != nil {
		panic("failed to auto migrate: " + err.Error())
	}

	DB = db
	return db
}
