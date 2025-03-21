//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	proc := MockProcess{}

	go func() {
		proc.Run() // Запуск процесса (блокирующий)
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)


	select {
	case <-sigint:
		fmt.Println("Попытка изящной остановки...")
		go proc.Stop()

		select {

		case <-time.After(3 * time.Second):
			fmt.Println("Процесс не остановился. Принудительное завершение.")
			os.Exit(1)
		case <-sigint:
			fmt.Println("Принудительное завершение.")
			os.Exit(1)
		}
	}
}