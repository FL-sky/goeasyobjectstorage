package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// 除了RabbitMQ提供的这个函数库，我们还自己实现了一个rabbitmq包，这是我们对“github.com/streadway/amqp”包的封装，
// 用于简化接口，代码见例2-11。

type RabbitMQ struct {
	channel  *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	// fmt.Println("New START")
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}

	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}

	q, e := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if e != nil {
		panic(e)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.conn = conn
	mq.Name = q.Name
	// fmt.Println("New END")
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	// fmt.Println("Bind START")
	e := q.channel.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if e != nil {
		panic(e)
	}
	q.exchange = exchange
	// fmt.Println("Bind END")
}

//New函数用于创建一个新的rabbitmq.RabbitMQ结构体，
// 该结构体的Bind方法可以将自己的消息队列和一个exchange绑定，所有发往该exchange的消息都能在自己的消息队列中被接收到。

func (q *RabbitMQ) Send(queue string, body interface{}) {
	// fmt.Println("Send START")
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
	// fmt.Println("Send END")
}

// Send方法可以往某个消息队列发送消息。
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	// fmt.Println("Publish START")
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
	// fmt.Println("Publish END")
}

// Publish方法可以往某个exchange发送消息。
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	// fmt.Println("Consume START")
	c, e := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	// fmt.Println("Consume END")
	return c
}

// Consume方法用于生成一个接收消息的go channel，使客户程序可以通过Go语言的原生机制接收队列中的消息。
func (q *RabbitMQ) Close() {
	// fmt.Println("Close START")
	q.channel.Close()
	q.conn.Close()
	// fmt.Println("Close END")
}

// Close方法用于关闭消息队列。
// 更多RabbitMQ接口的相关资料见其官方网站。
