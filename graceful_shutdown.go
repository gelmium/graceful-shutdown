package gfshutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Operation is a clean up function on shutting down.
type Operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it.
func GracefulShutdown(triggerCtx context.Context, timeout time.Duration, ops map[string]Operation) <-chan int {
	shutdownContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wait := make(chan int)
	go func() {
		s := make(chan os.Signal, 1)

		// these syscalls signal will trigger graceful shutdown.
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		// wait from both triggerCtx.Done and syscall s chanel.
		select {
		case <-triggerCtx.Done():
			log.Println("graceful shutdown is triggered within process")
		case signal := <-s:
			log.Println("received shutdown signal from system", signal)
		}

		log.Println("Shutting down...")

		// set timeout for the ops to be done to prevent system hang.
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			wait <- 1
			close(wait)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time.
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(shutdownContext); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		log.Printf("Graceful shutdown complete")
		wait <- 0
		close(wait)
	}()

	return wait
}
