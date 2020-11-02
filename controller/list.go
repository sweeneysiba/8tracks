package controller

import (
	"8tracks/config"
	"8tracks/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//PrettyPrint ...
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

//CreateList ...
func CreateList(c *gin.Context) {
	var list models.List
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	if list.Name == "" || list.Tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | name , tag for this playlist are mandatory"})
		return
	}
	timeStamp := time.Now().UnixNano()
	list.CreatedAt = timeStamp / int64(time.Millisecond)
	list.ID = timeStamp
	songsArray := []models.SongLists{}
	for _, val := range list.Songs {
		val.ListID = list.ID
		songDB, err := getSongInfo(val.SongID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		val.Name = songDB.Name
		val.Author = songDB.Author
		songsArray = append(songsArray, val)
	}
	list.Songs = songsArray
	PrettyPrint(list)

	if err := config.DB.Create(list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	// for _, val := range list.Songs {
	// 	val.ListID = list.ID
	// 	fmt.Println("inserting a song", val.ListID)
	// 	songDB, err := getSongInfo(val.SongID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 		return
	// 	}
	// 	val.Name = songDB.Name
	// 	val.Author = songDB.Author
	// 	PrettyPrint(val)
	// 	if err := config.DB.Create(val).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | Trying to add a song which doesn't exist '", "error": err})
	// 		return
	// 	}
	// 	fmt.Println("inserted a song", val.ListID)
	// }
	c.JSON(http.StatusOK, gin.H{"message": "List Created Successfully .", "status": 200})
	return

}

//DeleteList ...
func DeleteList(c *gin.Context) {
	listID := c.Param("id")
	fmt.Println("deleting list with id", listID)
	if listID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the listID in URl"})
		return
	}
	if err := config.DB.Where("list_id = ?", listID).Delete(&models.SongLists{}).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
			return
		}
	}
	if err := config.DB.Where("id = ?", listID).Delete(&models.List{}).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the list with id =" + listID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "List Deleted Successfully .", "status": 200})
	return
}

//GetList ...
func GetList(c *gin.Context) {
	listID := c.Param("id")
	var listItems []models.List
	query := "SELECT * FROM lists"
	if listID != "" {
		query += " where id ='" + listID + "'"
	}
	if err := config.DB.Raw(query).Scan(&listItems).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No List found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	for index, val := range listItems {
		songquery := "select * from song_lists where list_id= '" + strconv.Itoa(int(val.ID)) + "'"
		var songItems []models.SongLists
		if err := config.DB.Raw(songquery).Scan(&songItems).Error; err != nil {
			fmt.Println(err)
			if !gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
				return
			}
		}
		listItems[index].Songs = songItems
	}
	c.JSON(http.StatusOK, gin.H{"message": "list fetched Successfully", "Lists": listItems})
	return
}

//UpdateList ...
func UpdateList(c *gin.Context) {
	listID := c.Param("id")
	var list *models.UpdateList
	if listID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the listID in URl"})
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	var obj = make(map[string]interface{})
	if list.Name != "" {
		obj["name"] = list.Name
	}
	if list.Tag != "" {
		obj["tag"] = list.Tag
	}
	list.ID, _ = strconv.ParseInt(listID, 10, 64)
	PrettyPrint(list)
	if err := config.DB.Model(&models.List{}).Where("id = ?", listID).Updates(obj).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the list with id =" + listID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
			return
		}
	}
	for _, val := range list.Songs {
		songDB, err := getSongInfo(val.SongID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if val.Type == "insert" {
			fmt.Println("Inserting a record")

			var objInsert *models.SongLists
			objInsert = &models.SongLists{
				Name:   songDB.Name,
				Author: songDB.Author,
				ListID: list.ID,
				SongID: val.SongID,
			}
			PrettyPrint(objInsert)
			if err := config.DB.Create(objInsert).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
				return
			}
			fmt.Println("Inserted a record")

		} else if val.Type == "update" {
			fmt.Println("updating a record")
			objUpdate := make(map[string]interface{})
			if songDB.Name != "" {
				objUpdate["name"] = songDB.Name
			}
			if songDB.Author != "" {
				objUpdate["author"] = songDB.Author
			}
			updateSong := &models.SongList{
				ListID: val.ListID,
			}
			if err := config.DB.Model(&updateSong).Where("id = ?", val.ID).Updates(objUpdate).Error; err != nil {
				fmt.Println(err)
				if gorm.IsRecordNotFoundError(err) {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the list with id =" + listID, "error": err})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
					return
				}
			}
			fmt.Println("updated a record")

		} else if val.Type == "delete" {
			fmt.Println("deleting a record")

			if err := config.DB.Where("id = ?", val.ID).Delete(&models.SongList{}).Error; err != nil {
				fmt.Println(err)
				if gorm.IsRecordNotFoundError(err) {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the songs with id =" + strconv.Itoa(int(val.ID)), "error": err})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove songs from list ", "error": err})
					return
				}
			}
			fmt.Println("deleted a record")

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "list updated Successfully"})
	return
}

//ExploreLists ...
func ExploreLists(c *gin.Context) {
	lists := c.Param("tags")
	listArray := []string{}
	if lists != "" {
		listArray = strings.Split(lists, "+")
	}
	condition := ""
	if len(listArray) > 0 {
		condition = " where "
		for index, item := range listArray {
			if index != 0 {
				condition += " OR "
			}
			condition += " tag like '%" + item + "%'"
		}
	}
	var listItems []models.List
	query := "SELECT * FROM lists" + condition + " order by likes+plays desc "
	fmt.Println(query)
	if err := config.DB.Raw(query).Scan(&listItems).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No List found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	for index, val := range listItems {
		songquery := "select * from song_lists where list_id= '" + strconv.Itoa(int(val.ID)) + "'"
		var songItems []models.SongLists
		if err := config.DB.Raw(songquery).Scan(&songItems).Error; err != nil {
			fmt.Println(err)
			if !gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
				return
			}
		}
		listItems[index].Songs = songItems
	}
	var tags []models.List
	if err := config.DB.Raw("select distinct(tag) from lists").Scan(&tags).Error; err != nil {
		fmt.Println(err)
		if !gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
			return
		}
	}
	fmt.Println(tags)
	tagsArray := []string{}
	for _, val := range tags {
		tagsArray = append(tagsArray, strings.Split(val.Tag, "+")...)
	}
	tagsArray = removeDuplicate(tagsArray)
	c.JSON(http.StatusOK, gin.H{"message": "list fetched Successfully", "Lists": listItems, "Tags": tagsArray})
	return
}
func removeDuplicate(obj []string) []string {

	for i := 0; i < len(obj); i = i + 1 {
		for j := i + 1; j < len(obj); j = j + 1 {
			if reflect.DeepEqual(obj[i], obj[j]) {
				copy(obj[j:], obj[j+1:])
				obj[len(obj)-1] = ""
				obj = obj[:len(obj)-1]
			}
		}
	}

	return obj
}

func getSongInfo(ID int64) (models.Songs, error) {
	var songDb models.Songs
	query := "SELECT * FROM songs"
	if ID != 0 {
		query += " where id ='" + strconv.Itoa(int(ID)) + "'"
	}
	if err := config.DB.Raw(query).Scan(&songDb).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			return songDb, errors.New("updating a song that doesn't exist")
		}
		return songDb, errors.New("enable to fetch song details")
	}
	return songDb, nil
}

//AddPlay ...
func AddPlay(c *gin.Context) {
	listID := c.Param("id")
	if listID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the listID in URl"})
		return
	}
	var listItems models.List
	query := "SELECT * FROM lists"
	if listID != "" {
		query += " where id ='" + listID + "'"
	}
	if err := config.DB.Raw(query).Scan(&listItems).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No List found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	PrettyPrint(listItems)
	var obj = make(map[string]interface{})
	obj["plays"] = listItems.Plays + 1
	if err := config.DB.Model(&models.List{}).Where("id = ?", listID).Updates(obj).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the list with id =" + listID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to find list ", "error": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated the play number"})
	return
}

//UpdateLike ...
func UpdateLike(c *gin.Context) {
	listID := c.Param("id")

	if listID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please pass the listID in URl"})
		return
	}
	var data struct {
		IsLike string `json:"is_like"`
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	var listItems models.List
	query := "SELECT * FROM lists"
	if listID != "" {
		query += " where id ='" + listID + "'"
	}
	if err := config.DB.Raw(query).Scan(&listItems).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No List found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	PrettyPrint(listItems)
	var obj = make(map[string]interface{})
	if data.IsLike == "like" {
		obj["likes"] = listItems.Likes + 1
	} else {
		obj["likes"] = listItems.Likes - 1
	}
	if err := config.DB.Model(&models.List{}).Where("id = ?", listID).Updates(obj).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable find the list with id =" + listID, "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to find list ", "error": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated the like number"})
	return
}
