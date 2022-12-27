package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FormId   uint   //发送着
	TargetId int64  //接受者
	Type     int    //消息类型
	Media    int    //消息类型
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// todo http://www.jsons.cn/websocket
// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 需要，发送者i，接受者id，消息类型，发送内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并检验token等合法性
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isValida := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3 用户关系
	//4 userid跟node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5完成发送逻辑
	go sendProc(node)
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for true {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<", data)
	}
}

var uppsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	uppsendChan <- data
}
func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送携程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 31, 25),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-uppsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}

		}

	}
}

// 完成udp数据接收携程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		//case 2: //群发
		//case 3: //广播
		//case 4:
	}
}
func sendMsg(UserId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[UserId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
