package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gogazub/consumer/model"
)

type ICodeRunner interface {
	RunCode(cm model.CodeMessage) error
}

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

func (r CodeRunner) RunCode(cm model.CodeMessage) error {
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
		return err
	}
	fmt.Print("Container created\n")

	ctxIO, cancelIO := context.WithCancel(context.Background())
	defer cancelIO()
	fmt.Print("Attach to container\n")
	attach, err := r.cli.ContainerAttach(ctxIO, resp.ID, container.AttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attach.Close()

	fmt.Print("Write code to container\n")
	if _, err := io.WriteString(attach.Conn, cm.Code); err != nil {
		return fmt.Errorf("write stdin: %w", err)
	}
	if err := attach.CloseWrite(); err != nil {
		return fmt.Errorf("close stdin: %w", err)
	}

	fmt.Print("Start container\n")
	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	fmt.Print("Container started\n")

	go func() {
		_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader)
	}()
	fmt.Print("Container wait\n")
	_, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	err = <-errCh
	if err != nil {
		return err
	}

	fmt.Println("container finished:", resp.ID)
	return nil
}

func Test() {
	code := `
#include<iostream>

int main(){
	std::cout << "Hello, from container!";
	return 0;
}		
`

	cr, err := NewCodeRunner()
	fmt.Print("Code runner created\n")
	if err != nil {
		fmt.Printf("CodeRunner creation error: %s", err.Error())
		return
	}
	codemsg := model.CodeMessage{Code: code}
	fmt.Print("Run code...\n")
	err = cr.RunCode(codemsg)
	if err != nil {
		fmt.Printf("RunCode error: %s", err.Error())
		return
	}
}
