package main

import (
	"context"
	"io"
	"log"
	"os"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc/metadata"
)

func(c *ClientService) Upload(ctx context.Context, fName *string, platform *string) error {
	iFile, err := os.Open(*fName)
	if err != nil {
		return err
	}
	defer iFile.Close()

	ctx = metadata.AppendToOutgoingContext(ctx, "platform", *platform)
	
	stream, err := c.client.Upload(ctx)
	if err != nil {
		return err
	}

	buf := make([]byte, chunkSize)
	for {
		num, err := iFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := buf[:num]

		if err := stream.Send(&pb.FileUploadRequest{Chunk: chunk}); err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	c.setSession(res.Session)
	log.Printf("Successfully send file %s over gRPC, file stored as %s", *fName, res.Session )
	return nil
}