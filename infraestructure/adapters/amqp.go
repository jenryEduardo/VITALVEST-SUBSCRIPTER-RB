package adapters

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
	conn, err := amqp.Dial("amqp://admin:tu_password_muy_segura@44.205.97.30:5672/")
	failOnError(err, "No se pudo conectar a RabbitMQ")
	return &ConnAMQP{conn: conn}
}

func (r *ConnAMQP) ListenToQueue() {
	ch, err := r.conn.Channel()
	failOnError(err, "No se pudo abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"usuarios_cola",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "No se pudo declarar la cola")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
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

			// 1. BME280
			bmePayload := map[string]interface{}{
				"temperatura": sensorData.BME280.Temperatura,
				"presion":     sensorData.BME280.Presion,
				"humedad":     sensorData.BME280.Humedad,
			}
			sendPost("http://100.30.168.141/bme", bmePayload, "BME280")

			// 2. MPU6050
			mpuPayload := map[string]interface{}{
				"pasos": sensorData.MPU6050.Pasos,
			}
			sendPost("http://100.30.168.141/mpu", mpuPayload, "MPU6050")

			// 3. MLX90614
			mlxPayload := map[string]interface{}{
				"temperatura_ambiente": sensorData.MLX90614.TemperaturaAmbiente,
				"temp_objeto":          sensorData.MLX90614.TempObjeto,
			}
			sendPost("http://100.30.168.141/mlx", mlxPayload, "MLX90614")

				// 4. GSR
			gsrpayload := map[string]interface{}{
				"porcentaje":sensorData.GSR.Porcentaje,
			}
			sendPost("http://100.30.168.141/gsr",gsrpayload,"GSR")
		}
	}()

	<-forever
}

// Función auxiliar para enviar POST
func sendPost(url string, payload map[string]interface{}, tag string) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ [%s] Error al serializar JSON: %v", tag, err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("❌ [%s] Error al hacer POST: %v", tag, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("✅ [%s] Datos enviados. Status: %s", tag, resp.Status)
}
