package db

import (
	"fmt"

	"github.com/mohammadMghi/apiGolangGateway/models"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var db *gorm.DB
type Mysql struct{
 
}

func ConnectToDB() (*gorm.DB, error) {
	dsn := "root:852456@tcp(127.0.0.1:3306)/restapi?charset=utf8mb4&parseTime=True&loc=Local"
	
	if db != nil{
		return db, nil
	}
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    db.AutoMigrate(&models.Transaction{})
 
 

	if err != nil{
		return nil , err
	}

	return db , err


}
 

func Insert(model models.Transaction){
	db , err:= ConnectToDB()

 

	if err !=nil{
		fmt.Errorf(err.Error())
	}

 

	db.Save(model)
	
}