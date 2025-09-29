package producer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type SourceCode struct{ Code []byte }

type Result struct {
	StdOut   []byte
	StdErr   []byte
	ExitCode int
	TimedOut bool
}

type DockerConfig struct {
	Image   string
	RunLine string
	Timeout time.Duration // таймаут на docker run
}

type ImageManager struct {
	PullTimeout time.Duration
}

func NewImageManager() *ImageManager {
	return &ImageManager{PullTimeout: 20 * time.Minute}
}

func (im *ImageManager) Ensure(ctx context.Context, image string) error {
	// Есть локально?
	if err := exec.CommandContext(ctx, "docker", "image", "inspect", image).Run(); err == nil {
		return nil
	}

	pctx, cancel := context.WithTimeout(context.Background(), im.PullTimeout)
	defer cancel()

	cmd := exec.CommandContext(pctx, "docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Разрулим типичное: время истекло vs процесс убит (OOM/kill)
		switch {
		case errors.Is(pctx.Err(), context.DeadlineExceeded):
			return fmt.Errorf("docker pull timeout after %s", im.PullTimeout)
		default:
			return fmt.Errorf("docker pull %s failed: %w", image, err)
		}
	}
	return nil
}

type Runner struct {
	baseArgs     []string
	imageManager *ImageManager
}

func NewRunner(im *ImageManager) *Runner {
	return &Runner{
		imageManager: im,
		baseArgs: []string{
			"run", "--rm", "-i",
			"--pull", "never",
			"--network", "none",
			"--cpus", "1",
			"--memory", "256m",
			"--memory-swap", "256m",
			"--pids-limit", "64",
			"--security-opt", "no-new-privileges",
			"--cap-drop", "ALL",
			"--read-only",
			"--tmpfs", "/sandbox:rw,nosuid,nodev,mode=1777,size=64m",
			"--tmpfs", "/tmp:rw,nosuid,nodev,mode=1777,size=64m",
			"-w", "/sandbox",
			"--user", "65534:65534",
			"-e", "TMPDIR=/tmp",
			"-e", "HOME=/sandbox",
		},
	}
}

func (r *Runner) RunCodeInContainer(ctx context.Context, sc SourceCode, cfg DockerConfig) (Result, error) {
	// 1) гарантируем образ
	if err := r.imageManager.Ensure(ctx, cfg.Image); err != nil {
		return Result{}, fmt.Errorf("ensure image %q: %w", cfg.Image, err)
	}

	// 2) собираем docker args
	args := make([]string, 0, len(r.baseArgs)+4)
	args = append(args, r.baseArgs...)
	args = append(args, cfg.Image, "bash", "-lc", cfg.RunLine)

	// 3) таймаут на run (если не задан выше)
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = bytes.NewReader(sc.Code)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	res := Result{
		StdOut: stdout.Bytes(),
		StdErr: stderr.Bytes(),
	}

	if err != nil {
		// run убит по таймауту хоста
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			res.TimedOut = true
			res.ExitCode = -1
			return res, fmt.Errorf("host timeout exceeded: %w", err)
		}
		// docker вернул ненулевой код (например, 125  ошибка docker, 127  команда не найдена и т.д.)
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.ExitCode = ee.ExitCode()
			return res, nil
		}
		// иная системная ошибка
		return res, err
	}

	res.ExitCode = 0
	return res, nil
}

var cppCfg = DockerConfig{
	Image:   "gcc:14-bookworm",
	RunLine: `cat > main.cpp && g++ -O2 -pipe main.cpp -o prog && timeout 2s ./prog`,
	Timeout: 10 * time.Second,
}

func main() {
	im := NewImageManager()
	runner := NewRunner(im)

	code := []byte(`#include <iostream>
int main(){ std::cout << "Hello from C++\n"; return 0; }`)

	ctx := context.Background()
	res, err := runner.RunCodeInContainer(ctx, SourceCode{Code: code}, cppCfg)
	if err != nil {
		log.Printf("run error: %v", err)
	}
	log.Printf("exit=%d timeout=%v", res.ExitCode, res.TimedOut)
	log.Printf("STDOUT:\n%s", res.StdOut)
	log.Printf("STDERR:\n%s", res.StdErr)
}
