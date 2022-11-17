package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseConnection() *gorm.DB {
	dsn := "host=localhost user=postgres password=Kleopatra2003! dbname=account_balances_go_4 port=5433"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't connect to the database")
	}
	// defer db.Close()

	return db
}

func SetupDatabase(db *gorm.DB) {
	db.AutoMigrate(&User{}, &UserAccount{}, &CompanyAccount{}, &Service{}, &Order{}, &OrderItem{}, &Transaction{})
}
