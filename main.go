package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
) 

type User struct {
	ID		uint	`gorm:column:id;primaryKey`
	Name	string	`gorm:"column:name"`
	Email	string	`gorm:"column:email"`
	Age		string	`gorm:"column:age"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time	`gorm:"column:updatedAt"`
}

func main() {
	dsn := "root:@tcp(localhost:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	db.AutoMigrate(&User{})

	router := gin.Default()

	//Endpoint untuk menampilkan semua pengguna
	router.GET("/users", func(c *gin.Context) {
		var users []User
		db.Find(&users) //Mengambil semua data dari tabel users
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	//Endpoint untuk menampilkan pengguna berdasarkan id
	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

//Endpoint untuk menambahkan pengguna baru
router.POST("/users", func(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	db.Create(&user) //Menambahkan data ke database
	c.JSON(http.StatusCreated, gin.H{"data": user})
})

//Endpoint untuk mengedit pengguna berdasarkan id
router.PUT("/users/:id", func(c *gin.Context) {
	var user User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user) //Menyimpan perubahan ke database
	c.JSON(http.StatusOK, gin.H{"data": user})
})

//Endpoint untuk menghapus pengguna berdasarkan id
router.DELETE("/users/:id", func(c *gin.Context) {
	var user User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	db.Delete(&user) //Menghapus data dari databse
	c.JSON(http.StatusOK, gin.H{"data": "Pengguna Berhasil Dihapus"})
})

// Menjalankan server di port 3000
router.Run(":3000")

}