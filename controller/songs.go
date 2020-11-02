package controller

import (
	"8tracks/config"
	"8tracks/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//CreateSong ...
func CreateSong(c *gin.Context) {
	var Song models.Songs
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &Song)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	if Song.Name == "" || Song.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | name , Author is mandatory for adding a new song"})
		return
	}
	timeStamp := time.Now().UnixNano()
	Song.CreatedAt = timeStamp / int64(time.Millisecond)
	Song.ID = timeStamp
	// PrettyPrint(userInput)
	if err := config.DB.Create(Song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song Created Successfully .", "status": 200})
	return

}

//DeleteSong ...
func DeleteSong(c *gin.Context) {
	SongID := c.Param("id")
	fmt.Println("deleting Song with id", SongID)
	if SongID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the SongID in URl"})
		return
	}

	if err := config.DB.Where("song_id = ?", SongID).Delete(&models.SongLists{}).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
			return
		}
	}

	if err := config.DB.Where("id = ?", SongID).Delete(&models.Songs{}).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove song", "error": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song Deleted Successfully .", "status": 200})
	return
}

//GetSongList ...
func GetSongList(c *gin.Context) {
	SongID := c.Param("id")
	var SongList []models.Songs
	query := "SELECT * FROM songs"
	if SongID != "" {
		query += " where id ='" + SongID + "'"
	}
	if err := config.DB.Raw(query).Scan(&SongList).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No Songs found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Songs fetched Successfully", "Songs": SongList})
	return
}

//UpdateSong ...
func UpdateSong(c *gin.Context) {
	SongID := c.Param("id")
	var Songs *models.Songs
	if SongID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the SongID in URl"})
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &Songs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	var obj = make(map[string]interface{})
	if Songs.Name != "" {
		obj["name"] = Songs.Name
	}
	if Songs.Author != "" {
		obj["author"] = Songs.Author
	}
	Songs.ID, _ = strconv.ParseInt(SongID, 10, 64)
	PrettyPrint(Songs)
	if err := config.DB.Model(&models.Songs{}).Where("id = ?", SongID).Updates(obj).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the Song with id =" + SongID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from Songs ", "error": err})
			return
		}
	}

	if err := config.DB.Model(&models.SongLists{}).Where("song_id = ?", SongID).Updates(obj).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the Song with id =" + SongID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from Songs ", "error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated Successfully"})
	return
}
