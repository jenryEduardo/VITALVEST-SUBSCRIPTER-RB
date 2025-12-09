package main

import (
	"fmt"
	"net/http"
	"log"
	"sub/infraestructure/adapters"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! Your app is running on port 8082.")
}

func main() {
	// Inicializa la conexión y escucha de RabbitMQ en una goroutine
	// para que no bloquee el inicio del servidor HTTP.
	go func() {
		rabbit := adapters.NewConn()
		rabbit.ListenToQueue()
	}()

	// Configura y levanta el servidor HTTP en el puerto 8082
	http.HandleFunc("/", handler)
	port := ":8082"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	
	// log.Fatal hará que el programa termine si hay un error al iniciar el servidor
	log.Fatal(http.ListenAndServe(port, nil))
}