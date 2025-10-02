package processor

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gogazub/consumer/model"
)

type CodeProcessor struct {
}

func (cp *CodeProcessor) RunCode(code model.CodeMessage) {

}

func Test() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	res, err := cli.ImagePull(context.Background(), "hello-world", image.PullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	io.Copy(os.Stdout, res)
}
