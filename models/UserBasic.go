package models

import (
	"fmt"
	"ginChat/untils"
	"gorm.io/gorm"
	"time"
)

type GetUserLists struct {
	UserBasic  []*UserBasic
	UserNumber int
}

type UserBasic struct {
	UserId        string `json:"user_id" gorm:"colum:user_id"`
	Name          string `json:"name" gorm:"colum:name"`
	PassWord      string `json:"pass_word" gorm:"colum:pass_word"`
	Phone         string `json:"phone" gorm:"colum:phone" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `json:"email" gorm:"colum:email" valid:"email"`
	Identity      string `json:"identity" gorm:"colum:identity"`
	ClicentIp     string `json:"clicent_ip" gorm:"colum:clicent_ip"`
	LoginTime     string `json:"login_time" gorm:"colum:login_time"`
	HeartbeatTime string `json:"heartbeat_time" gorm:"colum:heartbeat_time"`
	LoginOutTime  string `json:"login_out_time" gorm:"colum:login_out_time"`
	IsLogout      bool   `json:"is_logout" gorm:"colum:is_logout"`
	DeviceInfo    string `json:"device_info" gorm:"colum:device_info"`
	Salt          string `json:"salt" gorm:"colum:salt"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

//CREATE TABLE `user_basic` (
//`user_id` varchar(25) NOT NULL COMMENT '用户ID',
//`name` varchar(45) NOT NULL COMMENT '用户名',
//`pass_word` varchar(45) NOT NULL COMMENT '密码',
//`phone` varchar(15) NOT NULL COMMENT '电话',
//`email` varchar(20) DEFAULT '' COMMENT '邮箱',
//`identity` varchar(40) DEFAULT '' COMMENT 'token',
//`clicent_ip` varchar(20) DEFAULT '' COMMENT '用户ip',
//`login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登陆时间',
//`heartbeat_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '更新时间',
//`login_out_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登陆时间',
//`is_logout` tinyint DEFAULT '1' COMMENT '是否退出 1:已下线 2:在线',
//`device_info` varchar(10) DEFAULT '' COMMENT '设备',
//`status` tinyint DEFAULT '1' COMMENT '用户状态 1:正常 2:删除',
//`salt` varchar(32) DEFAULT '' COMMENT '密码加盐',
//PRIMARY KEY (`user_id`),
//KEY `name` (`name`),
//KEY `phone` (`phone`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

func GetUserList() (getUserLists GetUserLists) {
	data := make([]*UserBasic, 10)
	untils.DB.Debug().Where("status = 1").Find(&data)
	getUserLists.UserBasic = data
	getUserLists.UserNumber = len(data)
	return getUserLists
}

func CreateUser(user UserBasic) *gorm.DB {
	return untils.DB.Debug().Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return untils.DB.Model(&user).Debug().Where("user_id = ?", user.UserId).Update("status", 2)
}

func GetUser(user UserBasic) UserBasic {
	var userBasic UserBasic
	untils.DB.Debug().Where("user_id = ?", user.UserId).
		Where("pass_word = ?", user.PassWord).
		Find(&userBasic)
	if userBasic.UserId != "" {
		//token加密
		token := fmt.Sprintf("%d", time.Now().Unix())
		identity := untils.Md5Encode(token)
		untils.DB.Model(&user).Debug().Where("user_id = ?", user.UserId).
			Update("identity", identity)
	}
	return userBasic
}
func GetSalt(userId string) string {
	var userBasic UserBasic
	untils.DB.Debug().Where("user_id = ?", userId).
		Find(&userBasic)
	return userBasic.Salt
}
