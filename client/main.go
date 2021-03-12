/* basic grpc client to test cloudflare timeouts
 */

package main

import (
	"context"
	"fmt"
	"time"

	testgrpc "grpc_test/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var localhost = false

// for cloudflare
var host = "grpc-test.matej.me:443"
var cert = "cf_root.pem"

func main() {
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		panic(err)
	}

	var conn *grpc.ClientConn
	if localhost {
		conn, err = grpc.Dial("127.0.0.1:2015",
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Second*3),
			grpc.WithBlock(),
		)
	} else {
		conn, err = grpc.Dial(host,
			grpc.WithTransportCredentials(creds),
			grpc.WithTimeout(time.Second*3),
			grpc.WithBlock(),
		)
	}

	if err != nil {
		panic(err)
	}

	client := testgrpc.NewTestServiceClient(conn)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*2))
	defer cancel()

	c, err := client.Bidirectional(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected", time.Now())
	ticker := time.NewTicker(time.Second * 3)

	for {
		// wait to send ping
		<-ticker.C

		fmt.Println("Sending ping")
		if err := c.Send(&testgrpc.Ping{Time: timestamppb.Now()}); err != nil {
			panic(err)
		}

		msg, err := c.Recv()
		if err != nil {
			panic(err)
		}

		fmt.Println("Received "+msg.Event.String(), time.Now())
	}
}
