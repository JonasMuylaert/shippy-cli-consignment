package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/JonasMuylaert/shippy/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(fn string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &consignment); err != nil {
		return nil, err
	}
	return consignment, nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("an error occured while connecting to the server: %v", err)
	}

	client := pb.NewShippingServiceClient(conn)

	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("couldn't parse file: %v", err)
	}

	res, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("failed creating consignment: %v", err)
	}

	log.Printf("Created consignment: %v", res.Created)

	all, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed getting consginments: %v", err)
	}
	for _, cons := range all.Consignments {
		log.Printf("Consignment: %v\n", cons)
	}
}
