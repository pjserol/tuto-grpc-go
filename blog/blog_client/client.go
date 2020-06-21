package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/pjserol/tuto-grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Client Start!")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	id, err := createBlog(c)
	if err != nil {
		log.Printf("Error createBlog: %v", err)
	}

	err = readBlog(c, "id123")
	if err != nil {
		log.Printf("Error readBlog: %v", err)
	}

	err = readBlog(c, id)
	if err != nil {
		log.Printf("Error readBlog: %v", err)
	}

	err = updateBlog(c, id)
	if err != nil {
		log.Printf("Error updateBlog: %v", err)
	}

	err = deleteBlog(c, id)
	if err != nil {
		log.Printf("Error deleteBlog: %v", err)
	}

	err = listBlog(c)
	if err != nil {
		log.Printf("Error listBlog: %v", err)
	}

	myPath := "/Users/pj/Downloads/"
	b, err := downloadImage(c, myPath+"kobe_logoo.jpeg")
	if err != nil {
		log.Printf("Error downloadImage: %v", err)
	}

	err = ioutil.WriteFile(myPath+"test.jpeg", b, 0644)
	if err != nil {
		log.Printf("Error WriteFile: %v", err)
	}
}

func createBlog(c blogpb.BlogServiceClient) (blogID string, err error) {
	log.Println("\n\n--Create blog--")

	blog := &blogpb.Blog{
		AuthorId: "Tic Tac",
		Title:    "My Blog",
		Content:  "Some content...",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		return "", err
	}

	fmt.Printf("Blog has been created: %v", res)

	return res.GetBlog().GetId(), nil
}

func readBlog(c blogpb.BlogServiceClient, id string) error {
	log.Println("\n\n--Reading blog--")

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: id})
	if err != nil {
		return err
	}

	fmt.Printf("Blog was read: %v \n", res)

	return nil
}

func updateBlog(c blogpb.BlogServiceClient, id string) error {
	log.Println("\n\n--Updating blog--")

	blog := &blogpb.Blog{
		Id:       id,
		AuthorId: "Ping Pong",
		Title:    "Blog PJ",
		Content:  "This is my first blog!",
	}

	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: blog})
	if err != nil {
		return err
	}

	fmt.Printf("Blog was updated: %v\n", res)

	return nil
}

func deleteBlog(c blogpb.BlogServiceClient, id string) error {
	log.Println("\n\n--Deleting blog--")
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: id})
	if err != nil {
		return err
	}

	fmt.Printf("Blog was deleted: %v \n", res)

	return nil
}

func listBlog(c blogpb.BlogServiceClient) error {
	log.Println("\n\n--List blog--")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println(res.GetBlog())
	}

	return nil
}

func downloadImage(c blogpb.BlogServiceClient, fileName string) ([]byte, error) {
	log.Println("\n\n--Dowanload image--")
	stream, err := c.DownloadImage(context.TODO(), &blogpb.DownloadImageRequest{
		FileName: fileName,
	})
	if err != nil {
		return nil, err
	}

	var img []byte
	for {
		chunkResponse, err := stream.Recv()
		if err == io.EOF {
			log.Println("received all chunks")
			break
		}
		if err != nil {
			log.Println("err receiving chunk:", err)
			break
		}

		img = chunkResponse.FileChunk
	}

	return img, nil
}
