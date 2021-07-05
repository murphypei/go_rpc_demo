package main

import (
	"bufio"
	"context"
	pb "demo/proto"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func request(srv pb.StudentService) {
	seq_name := [3]string{"java", "steven", "tony"}

	for i := 0; i < len(seq_name); i++ {
		rsp, err := srv.GetStudent(context.TODO(), &pb.StudentRequest{Name: seq_name[i]})
		if err != nil {
			log.Printf("Error in request: %v", err)
			return
		}
		log.Printf("request get name: %v, classes: %v, grade: %v", rsp.Name, rsp.Classes, rsp.Grade)
	}
}

func streamRequest(srv pb.StudentService) {
	stream, err := srv.StreamGetStudent(context.Background())
	if err != nil {
		log.Printf("Error in get stream: %v", err)
		return
	}

	// get request name from command line
	log.Printf("input request name")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "close" {
			break
		}

		if err := stream.Send(&pb.StudentRequest{Name: line}); err != nil {
			log.Printf("Error in stream send: %v", err)
		}

		if rsp, err := stream.Recv(); err != nil {
			log.Printf("Error in stream recv: %v", err)
			break
		} else {
			log.Printf("stream request get name: %v, classes: %v, grade: %v", rsp.Name, rsp.Classes, rsp.Grade)
		}

		if rsp, err := stream.Recv(); err != nil {
			log.Printf("Error in stream recv: %v", err)
			break
		} else {
			log.Printf("stream request get name: %v, classes: %v, grade: %v", rsp.Name, rsp.Classes, rsp.Grade)
		}
	}

	if err := stream.Close(); err != nil {
		log.Printf("Error in stream close: %v", err)
	}
}

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Name("student.client"),
		micro.Registry(consulReg),
	)
	service.Init()

	studentService := pb.NewStudentService("student_service", service.Client())
	request(studentService)
	streamRequest(studentService)

}
