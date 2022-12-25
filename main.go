package main

import (
	Router "ginChat/router"
	"ginChat/untils"
)

func main() {
	untils.InitConfig()
	untils.InitMysql()
	untils.InitRedis()
	r := Router.Router()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
