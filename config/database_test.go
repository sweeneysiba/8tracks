package config

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var err error

func TestStoreHotel(t *testing.T) {
	DB, err = gorm.Open("mysql", DbURL(BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer DB.Close()

}
