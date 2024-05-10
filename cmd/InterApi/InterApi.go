package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"database/sql"
	"encoding/json"

	config "inter/config"
	db "inter/dbpool"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"
)

const (
	name        = "InterApi"
	description = "포스 Api"
)

var dependencies = []string{"InterApi.service"}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: InterApi install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	proc()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case killSignal := <-interrupt:
			config.Stdlog.Println("Got signal:", killSignal)
			config.Stdlog.Println("Stoping DB Conntion : ", db.DB.Stats())
			config.Stdlog.Println("Stoping DB2 Conntion : ", db.DB2.Stats())
			defer db.DB.Close()
			defer db.DB2.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func main() {
	config.InitConfig()
	db.InitDB()
	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		config.Stdlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		config.Stdlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

type Mapper struct {
	MemId        string `json:"mem_id"`
	PosId        string `json:"pos_id"`
}

type Goods struct {
	GoodCd string
	GoodNm string
	SalePrc string
	HangPrc string
}

func proc(){
	config.Stdlog.Println("Inter API 시작")

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "테스트입니다.")
	})

	r.POST("/set_pos", func(c *gin.Context){
		db := db.DB2

		mp := Mapper{}
		ctx := c.Request.Context()
		err := c.ShouldBindJSON(&mp)
		if err != nil { 
			config.Stdlog.Println(err)
		}

		var cnt sql.NullInt64
		err = db.QueryRowContext(ctx, "select count(1) as cnt from key_mapper where mem_id = ?", mp.MemId).Scan(&cnt)
		if err != nil {
			config.Stdlog.Println(err)
		}
		if cnt.Int64 > 0 {
			_, err := db.ExecContext(ctx, "update key_mapper set pos_id = ? where mem_id = ?", mp.PosId, mp.MemId)
			if err != nil {
				config.Stdlog.Println(err)
				c.JSON(200, gin.H{
					"code": 0,
				})
			}
			c.JSON(200, gin.H{
				"code": 1,
			})
		} else {
			_, err := db.ExecContext(ctx, "insert into key_mapper values(?, ?)", mp.MemId, mp.PosId)
			if err != nil {
				config.Stdlog.Println(err)
				c.JSON(200, gin.H{
					"code": 0,
				})
			}
			c.JSON(200, gin.H{
				"code": 1,
			})
		}
	})
	r.POST("/get_goods", func(c *gin.Context){
		db := db.DB

		mp := Mapper{}
		ctx := c.Request.Context()
		err := c.ShouldBindJSON(&mp)
		if err != nil { 
			config.Stdlog.Println(err)
		}

		rows, err := db.QueryContext(ctx, "select cGoodcd, cGoodNm, fSalePrc, fHangPrc from GOOD1000LOG where cManID = ?", mp.PosId)
		if err != nil { 
			config.Stdlog.Println(err)
		}
		defer rows.Close()

		list := []Goods{}

		for rows.Next(){
			var goods Goods
			err := rows.Scan(&goods.GoodCd, &goods.GoodNm, &goods.SalePrc, &goods.HangPrc)
			if err != nil { 
				config.Stdlog.Println(err)
			}
			list = append(list, goods)
		}

		jsonBytes, err := json.Marshal(list)
		if err != nil {
			panic(err)
		}

		jsonString := string(jsonBytes)
		c.JSON(200, gin.H{
			"list": jsonString,
		})
	})

























	r.Run(":3333")
}