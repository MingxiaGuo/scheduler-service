package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scheduler-service/handlers"
	"scheduler-service/services"
	"syscall"
	"time"
	// "github.com/gin-gonic/gin"
)

func main() {

	port := flag.String("port", "8080", "Port for the HTTP server")
	bandwidth := flag.Int("bandwidth", 5, "Bandwidth of the scheduler")
	flag.Parse()

	taskService := services.NewTaskService(*bandwidth)
	taskHandler := handlers.NewTaskHandler(taskService)

	schedulerService := services.NewSchedulerService(taskService)
	schedulerService.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", taskHandler.SubmitTasks)
	mux.HandleFunc("/status", taskHandler.GetStatus)
	mux.HandleFunc("/scheduler", taskHandler.SwitchScheduler)

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on port %s with bandwidth %d\n", *port, *bandwidth)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	schedulerService.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exited")
}
