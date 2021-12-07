package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jsterling7/into-to-grpc/cachely"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// connect to the grpc server
	conn, err := grpc.Dial("127.0.0.1:5051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// create a new client
	c := cachely.NewCacheClient(conn)

	ctx := context.Background()

	// ------- exercise 1 -------

	// write a value
	_, err = c.Put(ctx, &cachely.PutRequest{
		Key:   "band",
		Value: []byte("The Beatles"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// #TODO: using the client, use the `Get` method to retrieve the value for the key `band`
	// # hint: remember to check the error
	getResponse, err := c.Get(ctx, &cachely.GetRequest{Key: "band"})
	code := status.Code(err)
	switch code {
	case codes.NotFound:
		log.Fatal("'band' was not found, put failed with err: ", err.Error())
	case codes.OK:
		fmt.Println("found value for 'band':")
		fmt.Println(getResponse.Value)
	default:
		log.Fatal(err)
	}

	// #TODO: delete the key for `band` using the `Delete` method from the client.
	_, err = c.Delete(ctx, &cachely.DeleteRequest{Key: "band"})
	code = status.Code(err)
	switch code {
	case codes.NotFound:
		log.Fatal("could not find 'band' key to delete with err: ", err.Error())
	case codes.OK:
		fmt.Println("delete request was successful")
	default:
		log.Fatal(err)
	}

	// check it was deleted
	getResponse, err = c.Get(ctx, &cachely.GetRequest{Key: "band"})
	code = status.Code(err)
	switch code {
	case codes.NotFound:
		fmt.Println("confirmed 'band' has been deleted")
	case codes.OK:
		log.Fatal("band code was found, delete failed with err: ", err.Error())
	default:
		log.Fatal(err)
	}

	// --------------

	// ------- exercise 2 -------
	// TODO call `PutMany` from the client and retrieve a reference to the `stream`

	//for i := 0; i < 100; i++ {
	//	key := fmt.Sprintf("key-%d", i)
	//	value := []byte(fmt.Sprintf("value-%d", i))
	//	// TODO: add the key/value with the stream
	//}

	// TODO: Close the stream

	/* Uncomment to verify your work
	// Verify all the values are there
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		item, err := c.Get(ctx, &cachely.GetRequest{
			Key: key,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received: key: %q, value: %q\n", key, string(item.GetValue()))
	}
	*/
}