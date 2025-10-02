package runner

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gogazub/consumer/model"
)

type CodeRunner struct {
}

func NewCodeProcessor() *CodeRunner {
	return &CodeRunner{}
}

func (cp *CodeRunner) RunCode(cm model.CodeMessage) {
	fmt.Printf("Mock: accepted code: %s to run", cm.Code)
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
