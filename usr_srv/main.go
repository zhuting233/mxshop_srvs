package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/usr_srv/global"
	"mxshop_srvs/usr_srv/handler"
	"mxshop_srvs/usr_srv/initialize"
	"mxshop_srvs/usr_srv/proto"
	"mxshop_srvs/usr_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	Ip := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	flag.Parse()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	if *Port == 0 {
		var err error
		if *Port, err = utils.GetFreePort(); err != nil {
			panic(err)
		}
	}
	zap.S().Infof("%s:%d\n", *Ip, *Port)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.57.1:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	//生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceId,
		Port:    *Port,
		Tags:    []string{"user", "srv", "grpc"},
		Address: "192.168.57.1",
		Check:   check,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Panic("注销失败")
	}
	zap.S().Info("注销成功")
}
