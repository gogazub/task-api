package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gogazub/consumer/model"
)

type CodeRunner struct {
	cli *client.Client
}

func NewCodeRunner() (*CodeRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &CodeRunner{cli: cli}, nil
}

func (r *CodeRunner) RunCode(cm model.CodeMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const image = "gcc:14-bookworm"

	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image:     image,
		Cmd:       []string{"bash", "-lc", "g++ -O2 -std=c++17 -x c++ - -o /tmp/a.out && /tmp/a.out"},
		Tty:       false,
		OpenStdin: true,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	attach, err := r.cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attach.Close()

	_, _ = attach.Conn.Write([]byte(cm.Code))
	attach.CloseWrite()

	if err := r.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	go func() {
		_, _ = io.Copy(os.Stdout, attach.Reader)
	}()

	_, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	err = <-errCh
	if err != nil {
		return err
	}

	fmt.Println("container finished:", resp.ID)
	return nil
}
