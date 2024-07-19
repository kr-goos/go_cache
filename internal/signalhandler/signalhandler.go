package signalhandler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignal() string {
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	defer signal.Stop(sigChan)

	sig := <-sigChan
	switch sig {
	case syscall.SIGHUP:
		return "interrupted by signal (SIGHUP)"
	case syscall.SIGINT:
		return "interrupted by signal (SIGINT)"
	case syscall.SIGTERM:
		return "interrupted by signal (SIGTERM)"
	case syscall.SIGQUIT:
		return "interrupted by signal (SIGQUIT)"
	}

	return fmt.Sprintf("interrupted by unknown signal : %d", sig)
}
