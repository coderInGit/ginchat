package service

import (
	"fmt"
	"ginChat/common"
	"ginChat/models"
	"ginChat/untils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	getUserLists := models.GetUserList()
	c.JSON(200, gin.H{
		"message": getUserLists,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/CreateUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	if err := common.Init("2021-12-03", 1); err != nil {
		log.Fatal("Init() failed, err = ", err)
		return
	}
	user.UserId = strconv.FormatInt(common.GenID(), 10)
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	user.Phone = c.Query("phone")
	user.Email = c.Query("email")
	user.LoginTime = time.Now().Format("2006-01-02 15:04:05")
	user.LoginOutTime = time.Now().Format("2006-01-02 15:04:05")
	user.HeartbeatTime = time.Now().Format("2006-01-02 15:04:05")
	user.Salt = fmt.Sprintf("%06d", rand.Int31())
	if repassword != password {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	user.PassWord = untils.MakePassword(password, user.Salt)
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"message": "新增用户成功",
	})
}

func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	user.UserId = c.Query("user_id")
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}

func Login(c *gin.Context) {
	user := models.UserBasic{}
	user.UserId = c.Query("user_id")
	passWord := c.Query("password")
	salt := models.GetSalt(user.UserId)
	if salt == "" {
		c.JSON(-1, gin.H{
			"message": "账号不存在",
		})
		return
	}
	user.PassWord = untils.MakePassword(passWord, salt)
	userBasic := models.GetUser(user)
	if userBasic.UserId == "" {
		c.JSON(-1, gin.H{
			"message": "账号或密码错误",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "登陆成功",
		})
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := untils.Subscribe(c, untils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
