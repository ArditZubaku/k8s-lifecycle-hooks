package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	println("Test Service started")

	time.Sleep(1 * time.Second)
	if data, err := os.ReadFile("/tmp/poststart.log"); err == nil {
		println("PostStart hook output:", string(data))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		println("Received termination signal, checking preStop hook output...")

		time.Sleep(2 * time.Second)

		if data, err := os.ReadFile("/tmp/prestop.log"); err == nil {
			println("PreStop hook output:", string(data))
		} else {
			println("Could not read preStop log:", err.Error())
		}

		println("Test Service stopped")
		os.Exit(0)
	}()

	time.Sleep(30 * time.Minute)
	println("Test Service finished normally")
}
