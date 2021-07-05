package main

import (
	"context"
	pb "demo/proto"
	"errors"
	"io"
	"log"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

// A derived class of StudentService, implements GetStudent and StreamGetStudent function.
type Handler struct {
	studentMap map[string]pb.Student
}

func NewHandler() *Handler {
	stuMap := map[string]pb.Student{
		"java":   pb.Student{Name: "java", Classes: "software", Grade: 80},
		"steven": pb.Student{Name: "steven", Classes: "computer", Grade: 70},
		"tony":   pb.Student{Name: "tony", Classes: "math", Grade: 60},
	}
	return &Handler{studentMap: stuMap}
}

func (handler *Handler) GetStudent(ctx context.Context, request *pb.StudentRequest, response *pb.Student) error {
	if request.Name == "" {
		return errors.New("request is empty")
	}

	if student, ok := handler.studentMap[request.Name]; ok {
		log.Printf("request name find: %v", request.Name)
		*response = student
	} else {
		log.Printf("request name not find: %v", request.Name)
		*response = pb.Student{Name: "None", Classes: "None", Grade: 0}
	}
	return nil
}

func (handler *Handler) StreamGetStudent(ctx context.Context, stream pb.StudentService_StreamGetStudentStream) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		response := pb.Student{Name: "None", Classes: "None", Grade: 0}
		if student, ok := handler.studentMap[request.Name]; ok {
			log.Printf("stream request name find: %v", request.Name)
			response = student
		} else {
			log.Printf("stream request name not find: %v", request.Name)

		}

		if err := stream.Send(&response); err != nil {
			return err
		}

		// add twice send
		time.Sleep(time.Duration(2) * time.Second)
		response = pb.Student{Name: "lazyman", Classes: "sleep", Grade: 100}
		if err := stream.Send(&response); err != nil {
			return err
		}
	}
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
	pb.RegisterStudentServiceHandler(service.Server(), NewHandler())

	err := service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
