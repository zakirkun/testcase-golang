package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User is a struct
type Artikel struct {
	ID       int `gorm:"primaryKey"`
	Nama     string
	Deskrpsi string
	Foto     string
	Tag      string
	Pembuat  string
}

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	var artikel Artikel

	db.AutoMigrate(&artikel)

	router := gin.Default()

	router.GET("/artikel", func(c *gin.Context) {
		var artikel []Artikel
		err := db.Find(&artikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, artikel)
	})

	router.GET("/artikel-search", func(c *gin.Context) {
		key := c.Query("key")

		var artikel []Artikel
		err := db.Where("tag like ?", "%"+key+"%").Or("nama like ?", "%"+key+"%").Or("pembuat like ?", "%"+key+"%").Find(&artikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, artikel)
	})

	router.GET("/artikel/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var artikel Artikel
		if err := c.ShouldBind(&artikel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		err = db.Where("id = ?", id).Find(&artikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if artikel.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
			return
		}

		c.JSON(http.StatusOK, artikel)
	})

	router.POST("/artikel", func(c *gin.Context) {
		var artikel Artikel
		if err := c.ShouldBind(&artikel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		inserArtikel := artikel
		inserArtikel.Nama = artikel.Nama
		inserArtikel.Deskrpsi = artikel.Deskrpsi
		inserArtikel.Foto = artikel.Foto
		inserArtikel.Tag = artikel.Tag
		inserArtikel.Pembuat = artikel.Pembuat

		err = db.Create(&inserArtikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "data was created",
		})
	})

	router.PUT("/artikel/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var a Artikel
		if err := c.ShouldBind(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var artikel Artikel
		err = db.Where("id = ?", id).First(&artikel).Error
		if err != nil {
			fmt.Println(err.Error())
		}

		if artikel.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
			return
		}

		artikel.Nama = a.Nama
		artikel.Deskrpsi = a.Deskrpsi
		artikel.Foto = a.Foto
		artikel.Tag = a.Tag
		artikel.Pembuat = a.Pembuat

		err = db.Save(&artikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Success updated",
		})
	})

	router.DELETE("/artikel/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var a Artikel
		if err := c.ShouldBind(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var artikel Artikel
		err = db.Where("id = ?", id).First(&artikel).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if artikel.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
			return
		}

		var art Artikel
		err = db.Delete(&art, id).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Success deleted",
		})

	})

	router.Run(":8080")
}
