package utils

import (
	"encoding/json"
	"os"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB Database `json:"db"`
}

type Database struct {
	Host        string   `json:"host"`
	Port        string   `json:"port"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Name        string   `json:"name"`
	Collections []string `json:"collections"`
}

var config Config

func GetConfig(configDir string) Config {
	if reflect.DeepEqual(config, Config{}) {
		file, err := os.Open(configDir)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			panic(err)
		}

		return config
	} else {
		return config
	}
}

func GetDB() *gorm.DB {
	dsn := config.DB.Username + ":" + config.DB.Password + "@tcp(" + config.DB.Host + ":" + config.DB.Port + ")/" + config.DB.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
