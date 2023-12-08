package main

import (
	"cmdb/models"
	"cmdb/routers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	err := models.InitMysql()
	if err != nil {
		panic(err)
	}
	defer models.Db.Close()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := routers.Routers()

	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 30 * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20

	srv := endless.NewServer(":18080", r)
	err = srv.ListenAndServe()
	if err != nil {
		if !strings.Contains(err.Error(), "Use of closed network connection") {
			log.Println(err)
			os.Exit(10)
		}
	}
}
