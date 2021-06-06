package main

import (
	"context"
	proto "demo/proto"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

type StudentManager struct {
}

func (mgr *StudentManager) GetStudent(ctx context.Context, request *proto.StudentRequest, response *proto.Student) error {
	studentMap := map[string]proto.Student{
		"java":   proto.Student{Name: "java", Classes: "software", Grade: 80},
		"steven": proto.Student{Name: "steven", Classes: "computer", Grade: 70},
		"tony":   proto.Student{Name: "tony", Classes: "math", Grade: 60},
	}
	if request.Name == "" {
		return errors.New("request is empty")
	}

	student := studentMap[request.Name]
	if student.Name != "" {
		fmt.Println(student.Name, student.Classes, student.Grade)
		*response = student
		return nil
	}

	return errors.New("not find student")
}

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Name("student_service"),
		micro.Version("v1.0.1"),
		micro.Registry(consulReg),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
	)
	service.Init()
	proto.RegisterStudentServiceHandler(service.Server(), new(StudentManager))

	err := service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
