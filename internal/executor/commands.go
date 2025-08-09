package executor

import (
	"bufio"
	"os/exec"
)

func ExecuteWithStreaming(comd string, stdoutCh, stderrCh chan string, doneCh chan error) {
	cmd := exec.Command("sh", "-c", comd) // use shell for more complex commands

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		doneCh <- err
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		doneCh <- err
		return
	}

	if err := cmd.Start(); err != nil {
		doneCh <- err
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdoutCh <- scanner.Text()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrCh <- scanner.Text()
		}
	}()

	err = cmd.Wait()
	doneCh <- err
}
