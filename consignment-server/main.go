//go:generate protoc -I=proto/consignment-service --go_out=plugins=grpc:./proto/dest/consignment-service proto/consignment-service/consignment.proto

package main

import (
	"context"
	"log"
	"net"

	pb "github.com/cuizw911/Shippy/consignment-server/proto/dest/consignment-service"
	"google.golang.org/grpc"
)

const Addr = ":50051"

// IRepository ... 仓库接口
type IRepository interface {
	Create(cons *pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment // 获取仓库中所有的货物
}

// Repository ... 仓库结构体，实现了IRepository接口
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(c *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, c)
	return c, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// 定义微服务
type service struct {
	repo Repository
}

// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	all := s.repo.GetAll()
	resp := &pb.Response{Consignments: all}
	return resp, nil
}

func main() {
	listener, e := net.Listen("tcp", Addr)
	if e != nil {
		log.Fatalf("failed to listen: %v", e)
	}

	log.Printf("listen on: %s \n", Addr)

	server := grpc.NewServer()
	repo := Repository{}

	pb.RegisterShippingServiceServer(server, &service{repo})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
