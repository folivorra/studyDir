package main

import (
	"context"
	"fmt"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/controller"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	store := storage.NewInMemoryStorage()                 // хранилище
	itemController := controller.NewItemController(store) // контроллер
	router := mux.NewRouter()                             // маршрутизатор
	itemController.RegisterRoutes(router)                 // регистрация маршрутов по заданным методам

	srv := &http.Server{ // создаем объект сервера
		Addr:    ":8080", // который будет слушать 8080 порт
		Handler: router,  // и обрабатываться маршрутизатором router
	}

	go func() {
		fmt.Println("Listening on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// запускаем сервер в горутине чтобы не заблокироваться в main
	// сервер слушает порт и при возникновении ошибки (кроме ошибки graceful shutdown) аварийно завершает работу

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Gracefully shutdown the server...")
	// создаем канал прерывания, чтобы корректно обрабатывать нажатие Ctrl+C, вызывая GS

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// задаем контекст отмены в 5 секунд, откладывая освобождение ресурсов

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	// даем серверу мягко завершится за эти 5 секунд, иначе завершаем аварийно
	log.Println("Server exiting")
}
