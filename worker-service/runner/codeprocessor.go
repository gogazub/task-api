package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gogazub/consumer/model"
)

type ICodeRunner interface {
	RunCode(cm model.CodeMessage) model.Result 
}

// CodeRunner executes code in docker container. Redirects stderr and stdout from container
type CodeRunner struct {
	cli *client.Client
}

// NewCodeRunner create new CodeRunner
func NewCodeRunner() (*CodeRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &CodeRunner{cli: cli}, nil
}

// RunCode run code in container and return model.Result
func (r CodeRunner) RunCode(cm model.CodeMessage) model.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const image = "gcc:14-bookworm"
	fmt.Print("Create container\n")
	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image:     image,
		Cmd:       []string{"bash", "-lc", "g++ -O2 -std=c++17 -x c++ - -o /tmp/a.out && /tmp/a.out"},
		Tty:       false,
		OpenStdin: true,
		StdinOnce: true,
	}, nil, nil, nil, "")
	if err != nil {
		return model.Result{Error: []byte(err.Error())}
	}
	fmt.Print("Container created\n")

	ctxIO, cancelIO := context.WithCancel(context.Background())
	defer cancelIO()
	fmt.Print("Attach to container\n")
	attach, err := r.cli.ContainerAttach(ctxIO, resp.ID, container.AttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		return model.Result{Error: []byte(err.Error())}
	}
	defer attach.Close()

	fmt.Print("Write code to container\n")
	if _, err := io.WriteString(attach.Conn, cm.Code); err != nil {
		return model.Result{Error: []byte(fmt.Sprintf("write stdin: %v", err))}
	}
	if err := attach.CloseWrite(); err != nil {
		return model.Result{Error: []byte(fmt.Sprintf("close stdin: %v", err))}
	}

	fmt.Print("Start container\n")
	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return model.Result{Error: []byte(err.Error())}
	}
	fmt.Print("Container started\n")

	var stdoutBuf, stderrBuf bytes.Buffer

	fmt.Print("Copy container output\n")
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attach.Reader)
	if err != nil {
		return model.Result{Error: []byte(fmt.Sprintf("copy output: %v", err))}
	}

	fmt.Print("Container wait\n")
	_, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	err = <-errCh
	if err != nil {
		return model.Result{
			Id : cm.Id,
			Error:  []byte(err.Error()),
			Output: stdoutBuf.Bytes(),
		}
	}

	fmt.Println("container finished:", resp.ID)
	
	return model.Result{
		Output: stdoutBuf.Bytes(),
		Error:  stderrBuf.Bytes(),
	}
}


