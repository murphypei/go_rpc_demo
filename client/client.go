package main

import (
	"context"
	proto "demo/proto"
	"fmt"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Name("student.client"),
		micro.Registry(consulReg),
	)
	service.Init()

	studentService := proto.NewStudentService("student_service", service.Client())
	res, err := studentService.GetStudent(context.TODO(), &proto.StudentRequest{Name: "java"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Name)
	fmt.Println(res.Classes)
	fmt.Println(res.Grade)
	time.Sleep(10 * time.Second)
}
