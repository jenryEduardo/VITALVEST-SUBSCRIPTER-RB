package adapters

import (
	"encoding/json"
	"log"
	"sub/domain"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type ConnAMQP struct {
	conn *amqp.Connection
}

func NewConn() *ConnAMQP {
	conn, err := amqp.Dial("amqp://guest:guest@13.220.254.127:5672/")
	failOnError(err, "No se pudo conectar a RabbitMQ")
	return &ConnAMQP{conn: conn}
}

func (r *ConnAMQP) ListenToQueue() {
	ch, err := r.conn.Channel()
	failOnError(err, "No se pudo abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"usuarios_cola", // nombre de la cola
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // argumentos
	)
	failOnError(err, "No se pudo declarar la cola")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-ack (puedes poner false si quieres controlar tú la confirmación)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "No se pudo registrar el consumidor")

	log.Println("📥 Esperando mensajes...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var sensorData domain.DatosSensor
			err := json.Unmarshal(d.Body, &sensorData)
			if err != nil {
				log.Printf("❌ Error al deserializar mensaje: %v", err)
				continue
			}

			log.Printf("📩 Recibido: %+v", sensorData)

			// Aquí puedes hacer algo con `sensorData`, por ejemplo guardarlo en DB


			

		}
	}()

	<-forever // Bloquea el main
}
