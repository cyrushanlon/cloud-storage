package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type configuration struct {
	Project  string `json:"project"`
	SyncPath string `json:"syncpath"`
	Bucket   string `json:"bucket"`
}

func loadConfig() configuration {
	var out configuration
	file, err := os.Open("config")
	if err != nil {
		panic(err) //we don't care if the config failed to load
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err) //we don't care if the config failed to load
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		panic(err) //we don't care if the config failed to load
	}

	return out
}

type bucketHandler struct {
	ctx    context.Context
	name   string
	handle *storage.BucketHandle
}

func createHandler(bucketName string) (bucketHandler, error) {

	out := bucketHandler{}
	out.ctx = context.Background()

	client, err := storage.NewClient(out.ctx)
	if err != nil {
		return out, err
	}

	out.name = bucketName
	out.handle = client.Bucket(out.name)

	return out, nil
}

func (b *bucketHandler) getObjects() error {

	objs := b.handle.Objects(b.ctx, nil)

	for {
		attrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		log.Println(attrs.Name)
	}

	return nil
}

func main() {
	config := loadConfig()

	log.Println("config:", config)

	b, err := createHandler(config.Bucket)
	if err != nil {
		panic(err)
	}

	err = b.getObjects()
	if err != nil {
		panic(err)
	}
}
