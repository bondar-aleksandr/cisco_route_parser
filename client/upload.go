package main

import (
	"context"
	"io"
	"log"
	"os"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
)

func Upload(ctx context.Context, c pb.RouteParserClient, fName *string, platform *string) (*string, error) {
	iFile, err := os.Open(*fName)
	if err != nil {
		return nil, err
	}
	defer iFile.Close()
	
	stream, err := c.FileTransfer(ctx)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	for {
		num, err := iFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		chunk := buf[:num]

		if err := stream.Send(&pb.FileUploadRequest{Platform: *platform, Chunk: chunk}); err != nil {
			return nil, err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully send file %s over gRPC, file stored as %s", *fName, res.FileName )
	return &res.FileName, nil
}