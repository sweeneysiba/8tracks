//main.go
package main

import (
	"8tracks/config"
	"8tracks/models"
	Routes "8tracks/routes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var err error

func Sum(x int, y int) int {
	return x + y
}

func main() {
	config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer config.DB.Close()
	config.DB.AutoMigrate(&models.List{}, &models.SongList{}, &models.Songs{})
	config.DB.Model(&models.SongList{}).AddForeignKey("list_id", "lists(id)", "CASCADE", "CASCADE")
	config.DB.Model(&models.SongList{}).AddForeignKey("song_id", "songs(id)", "CASCADE", "CASCADE")
	r := Routes.SetupRouter()
	r.Run(":8080")
}
