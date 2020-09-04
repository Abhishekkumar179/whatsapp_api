// // package main

// // import (
// // 	"fmt"
// // 	"log"
// // 	"net/http"

// // 	"github.com/gorilla/websocket"
// // )

// // var upgrader = websocket.Upgrader{
// // 	ReadBufferSize:  1024,
// // 	WriteBufferSize: 1024,
// // }

// // func homePage(w http.ResponseWriter, r *http.Request) {
// // 	fmt.Fprintf(w, "Home Page")
// // }
// // func setupRoutes() {
// // 	http.HandleFunc("/", homePage)
// // 	http.HandleFunc("/ws", wsEndpoint)
// // }

// // func reader(conn *websocket.Conn) {
// // 	for {
// // 		// read in a message
// // 		messageType, p, err := conn.ReadMessage()
// // 		if err != nil {
// // 			log.Println(err)
// // 			return
// // 		}
// // 		// print out that message for clarity
// // 		fmt.Println(string(p))

// // 		if err := conn.WriteMessage(messageType, p); err != nil {
// // 			log.Println(err)
// // 			return
// // 		}

// // 	}
// // }
// // func wsEndpoint(w http.ResponseWriter, r *http.Request) {
// // 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

// // 	// upgrade this connection to a WebSocket
// // 	// connection
// // 	ws, err := upgrader.Upgrade(w, r, nil)
// // 	if err != nil {
// // 		log.Println(err)
// // 	}
// // 	fmt.Println("connected")
// // 	reader(ws)
// // }
// // func main() {
// // 	fmt.Println("Hello World")
// // 	setupRoutes()
// // 	log.Fatal(http.ListenAndServe(":8080", nil))
// // }

// // // d23095efa00166fff3af155a
// // // 19709560f99d553e86939ea6
// // // e3c249966d31f73b943fd168
// // // 93b6e7aac00d360dec27cc1b
// // // 866f9006f016e19489194921
// // // 454b8b4a14d40e44326875bb
// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/streadway/amqp"
// )

// func main() {
// 	// url := os.Getenv("CLOUDAMQP_URL")
// 	// if url == "" {
// 	// 	url = "amqp://localhost"
// 	// }
// 	connection, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	defer connection.Close()

// 	//go func(connection *amqp.Connection) {

// 	channel, _ := connection.Channel()
// 	defer channel.Close()
// 	errs := channel.ExchangeDeclare(
// 		"Exchange", // name
// 		"direct",   // type
// 		true,       // durable
// 		false,      // auto-deleted
// 		false,      // internal
// 		false,      // no-wait
// 		nil,        // arguments
// 	)
// 	if errs != nil {
// 		fmt.Println("exchange not created")
// 	}
// 	durable, exclusive := true, false
// 	autoDelete, noWait := false, false
// 	q, _ := channel.QueueDeclare("test", durable, autoDelete, exclusive, noWait, nil)
// 	fmt.Println(q.Consumers, q.Messages)
// 	errq := channel.Qos(q.Messages, 0, false)
// 	if errq != nil {
// 		log.Fatalf("basic.qos: %v", errq)
// 	}
// 	err := channel.QueueBind(q.Name, "", "Exchange", false, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	//autoAck, exclusive, noLocal, noWait := false, false, false, false
// 	// messages, _, err := channel.Get("test", false)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	pages, err := channel.Consume("test", "", false, false, false, false, nil)
// 	if err != nil {
// 		log.Fatalf("basic.consume: %v", err)
// 	}

// 	go func() {
// 		for log := range pages {
// 			//fmt.Println("messages ", string(messages.Body))
// 			fmt.Println("log ", string(log.Body))
// 			fmt.Println("rede = ", log.Redelivered, " exc =", log.Exchange, " msgcount = ", log.MessageCount, " msgid = ", log.MessageId, " msgReplyto =", log.ReplyTo)

// 			// ... this consumer is responsible for sending pages per log
// 			//log.Acknowledger.Ack(23,true)
// 			log.Ack(true)
// 		}
// 	}()
// 	// err = channel.Qos(1, 0, false)
// 	// if err != nil {
// 	// 	log.Fatalf("basic.qos: %v", err)
// 	// }

// 	//messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
// 	//multiAck := false
// 	// for msg := range messages {
// 	// 	fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
// 	// 	msg.Ack(multiAck)
// 	// }
// 	//}(connection)

// 	// go func(con *amqp.Connection) {
// 	// 	timer := time.NewTicker(1 * time.Second)
// 	// 	channel, _ := connection.Channel()

// 	// 	for t := range timer.C {
// 	// 		msg := amqp.Publishing{
// 	// 			DeliveryMode: 1,
// 	// 			Timestamp:    t,
// 	// 			ContentType:  "text/plain",
// 	// 			Body:         []byte("Hello world"),
// 	// 		}
// 	// 		mandatory, immediate := false, false
// 	// 		channel.Publish("amq.topic", "ping", mandatory, immediate, msg)
// 	// 	}
// 	// }(connection)

// 	//select {}
// }
package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
