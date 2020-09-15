package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	pb "github.com/cuizw911/Shippy/consignment-server/proto/dest/consignment-service"
	"google.golang.org/grpc"
)

const (
	ADDRESS         = "localhost:50051"
	DefaultInfoFile = "consignment-client/consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(bytes, &consignment)
	if err != nil {
		log.Println(err)
		return nil, errors.New("unmarshal consignment.json file error")
	}
	return consignment, nil
}

func main() {

	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to grpc error: %v", err)
	}

	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)

	// 解析json信息
	consignment, err := parseFile(DefaultInfoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 调用GRPC
	response, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	log.Printf("created: %t", response.Created)

	ret, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("get consignments error: %v", err)
	}

	log.Printf("consignments: %v", ret.Consignments)
}
