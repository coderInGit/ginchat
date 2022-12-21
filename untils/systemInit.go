package untils

import (
	"fmt"
	"ginChat/common"
	"ginChat/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func InitMysql() {
	DB, err := gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	userBasic := &models.UserBasic{}
	userBasic.Name = "申专"
	if err := common.Init("2021-12-03", 1); err != nil {
		fmt.Println("Init() failed, err = ", err)
		return
	}
	userBasic.UserId = common.GenID()
	DB.Debug().Select("user_id", "name").Omit("created_at", "updated_at").Create(userBasic)
	os.Exit(22)
}
