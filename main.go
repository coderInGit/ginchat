package main

import (
	"ginChat/router"
	"ginChat/untils"
)

func main() {
	untils.InitConfig()
	untils.InitMysql()
	r := Router.Router()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
