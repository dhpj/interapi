package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
			defer db.DB.Close()
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

func proc(){
	config.Stdlog.Println("Inter API 시작")

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "테스트입니다.")
	})
	r.Run(":3333")
}