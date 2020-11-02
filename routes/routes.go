package routes

import (
	"8tracks/controller"

	"github.com/gin-gonic/gin"
)

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	curdGrpList := r.Group("/list")
	{
		curdGrpList.POST("create", func(c *gin.Context) {
			controller.CreateList(c)
		})
		curdGrpList.DELETE("delete/:id", func(c *gin.Context) {
			controller.DeleteList(c)
		})
		curdGrpList.GET("get/:id", func(c *gin.Context) {
			controller.GetList(c)
		})
		curdGrpList.GET("/", func(c *gin.Context) {
			controller.GetList(c)
		})
		curdGrpList.PUT("update/:id", func(c *gin.Context) {
			controller.UpdateList(c)
		})

		curdGrpList.PATCH("updateLike/:id", func(c *gin.Context) {
			controller.UpdateLike(c)
		})
		curdGrpList.PATCH("addPlay/:id", func(c *gin.Context) {
			controller.AddPlay(c)
		})
	}
	explore := r.Group("/explore")
	{
		explore.GET("/:tags", func(c *gin.Context) {
			controller.ExploreLists(c)
		})
		explore.GET("", func(c *gin.Context) {
			controller.ExploreLists(c)
		})
	}
	curdGrpSong := r.Group("/songs")
	{
		curdGrpSong.POST("create", func(c *gin.Context) {
			controller.CreateSong(c)
		})
		curdGrpSong.DELETE("delete/:id", func(c *gin.Context) {
			controller.DeleteSong(c)
		})
		curdGrpSong.GET("get/:id", func(c *gin.Context) {
			controller.GetSongList(c)
		})
		curdGrpSong.GET("/", func(c *gin.Context) {
			controller.GetSongList(c)
		})
		curdGrpSong.PUT("update/:id", func(c *gin.Context) {
			controller.UpdateSong(c)
		})
	}
	return r
}
