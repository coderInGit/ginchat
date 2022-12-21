package models

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	UserId        int64  `json:"user_id" gorm:"not null"`
	Name          string `json:"name" gorm:"not null"`
	PassWord      string `json:"pass_word" gorm:"not null"`
	Phone         string `json:"phone" gorm:"not null"`
	Email         string `json:"email" gorm:"not null"`
	Identity      string `json:"identity" gorm:"not null"`
	ClicentIp     string `json:"clicent_ip" gorm:"not null"`
	LoginTime     uint64 `json:"login_time" gorm:"not null"`
	HeartbeatTime uint64 `json:"heartbeat_time" gorm:"not null"`
	LoginOutTime  uint64 `json:"login_out_time" gorm:"not null"`
	IsLogout      bool   `json:"is_logout" gorm:"not null"`
	DeviceInfo    string `json:"device_info" gorm:"not null"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

//CREATE TABLE `user_basic` (
//`user_id` BIGINT NOT NULL  comment "用户ID",
//`name` VARCHAR(45) NOT NULL comment "用户名",
//`pass_word` VARCHAR(45)  NULL comment "密码",
//`phone` TINYINT  NULL comment "电话",
//`email` VARCHAR(20)  NULL DEFAULT "" comment "邮箱",
//`identity` TINYINT  NULL DEFAULT 1 comment "身份",
//`clicent_ip` VARCHAR(20)  NULL DEFAULT "" comment "用户ip",
//`login_time` DateTime DEFAULT CURRENT_TIMESTAMP   comment "登陆时间",
//`heartbeat_time` DateTime DEFAULT CURRENT_TIMESTAMP comment "",
//`login_out_time` DateTime DEFAULT CURRENT_TIMESTAMP   comment "登陆时间",
//`is_logout` TINYINT  NULL DEFAULT 1 comment "是否退出 1:已下线 2:在线",
//`device_info` VARCHAR(10)  NULL DEFAULT "" comment "设备",
//PRIMARY KEY (`user_id`),
//KEY `name` (`name`),
//KEY `phone` (`phone`)
//)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment="用户表";

//func GetUserList() []*UserBasic {
//	data := make([]*UserBasic, 10)
//	untils.DB.Find(&data)
//	return data
//}
