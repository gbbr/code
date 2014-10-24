// This sample program demonstrates how to use a channel to
// monitor the amount of time the program is running and terminate
// the program if it runs too long.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// shutdown provides system wide notification.
var shutdown = make(chan bool)

// workChan will retrieve any pending errors
var workChan = make(chan error)

var numWorkers = 10

// Timer runs

// main is the entry point for all Go programs.
func main() {
	// Launch the process.
	fmt.Println("Launching Processors")
	go startWork()

	err := ControlLoop()
	if err != nil {
		log.Printf("Process ended with message: %s", err)
	}

	// Program finished.
	fmt.Println("Process Ended")
}

// startWork provides the main program logic for the program.
func startWork() {
	log.Println("Processor - Starting")

	// Perform the work.
	workChan <- funnelWorkers()

	log.Println("Processor - Ended")
	// Capture any potential panic.
	if r := recover(); r != nil {
		log.Println("Processor - Panic", r)
	}
}

func ControlLoop() error {
	// sigChan receives os signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			// Interrupt event signaled by the operation system.
			log.Println("Interrupted - shutting down.")
			shutdown <- true

			// Set the channel to nil so we no longer process
			// any more of these events.
			sigChan = nil

		case <-time.After(5 * time.Second):
			// We have taken too much time.
			log.Println("Timeout - Exiting")
			os.Exit(1)

		case err := <-workChan:
			return err
		}
	}
}

// gotShutdown checks the stop flag to determine
// if we have been asked to interrupt processing.
func gotShutdown() bool {
	select {
	case <-shutdown:
		// We have been asked to stop cleanly.
		log.Println("Received stop signal.")
		return true

	default:
		return false
	}
}

// doWork simulates task work.
func funnelWorkers() error {
	time.Sleep(2 * time.Second)
	log.Println("Finished Task 1")

	if gotShutdown() {
		return errors.New("Early Shutdown")
	}

	time.Sleep(1 * time.Second)
	log.Println("Finished Task 2")

	if gotShutdown() {
		return errors.New("Early Shutdown")
	}

	time.Sleep(1 * time.Second)
	log.Println("Processor - Task 3")

	return nil
}
