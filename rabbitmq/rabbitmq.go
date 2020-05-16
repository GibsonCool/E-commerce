package _0_RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url 格式  amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const (
	MQURL           = "amqp://double:double@127.0.0.1:5672/test-double"
	SimpleQueueName = "doubleSimple"

	PubSubExchangeName = "newProduct"

	RoutingExchangeName = "exRouting"
	RoutingKey1         = "ex_one"
	RoutingKey2         = "ex_two"

	TopicExchangeName = "exTopic"
	TopicKey1         = "double.topic.one"
	TopicKey2         = "double.topic.two"
)

var (
	Done        chan bool
	DonePubSub  chan bool
	DoneRouting chan bool
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// binding key
	Key string
	// 连接信息
	Mqurl string
}

// 创建 RabbitMQ 实例
func NewRabbitMq(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	// 创建 rabbitmq 链接
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "获取 channel 失败")
	return rabbitmq
}

// 断开 channel/connection
func (r *RabbitMQ) Destory() {
	_ = r.conn.Close()
	_ = r.channel.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Printf("%s:%v", message, err)
		panic(fmt.Sprintf("%s:%v", message, err))
	}
}

// 简单模式 Step1：创建简单模式下 RabbitMQ 实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//简单模式交换机使用的是默认的
	return NewRabbitMq(queueName, "", "")
}

// 简单模式 Step2：简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在则会自动创建，如果存在则跳过创建
	// 好处：保证队列存在，消息能发送到队列中
	_, e := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否会自动删除：当最后一个消费者断开后，是否将队列中的消息清除
		false,
		/*
			是否具有排他性，意思只有自己可见，其他人不可用
			RabbitMQ:排他性队列（Exclusive Queue) : https://www.cnblogs.com/rader/archive/2012/06/28/2567779.html
		*/
		false,
		// 是否阻塞
		false,
		// 其他额外属性
		nil,
	)
	if e != nil {
		fmt.Println(e)
	}

	// 2.发送消息到队列中
	e = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果是 true，根据 exchange 类型 routkey 规则，如果无法找得到符合条件的队列那么会把消息返回给发送者
		false,
		// 如果是 true，当 exchange 发送消息到队列后发现队列上没有绑定消费者，会将消息返送回给发送者
		false,
		// 真正的消息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if e != nil {
		fmt.Println(e)
	}

}

// 简单模式 Step3：简单模式下消费代码
func (r *RabbitMQ) ConsumeSimple() {
	// 无论生产还是消费，第一步都是尝试先申请队列
	_, e := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否会自动删除：当最后一个消费者断开后，是否将队列中的消息清除
		false,
		// 是否具有排他性，意思只有自己可见，其他人不可用
		false,
		// 是否阻塞
		false,
		// 其他额外属性
		nil,
	)
	if e != nil {
		fmt.Println(e)
	}

	Done = make(chan bool)
	msgs, e := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者，消费者处理器名称
		"",
		//是否自动应答通知已收到消息
		true,
		//是否排他性,非唯一的消费者，其他消费者处理器也可以去竞争这个队列里面的消息任务
		false,
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,
		//是否阻塞
		false,
		nil,
	)

	go func() {
		for d := range msgs {
			// 这里接收到消息，实现处理逻辑
			log.Printf("接收到消息：%s", d.Body)
		}
	}()

	log.Printf("消费者已开启，等待消息产生。。。")
	<-Done

	r.Destory()
	log.Printf("消费者关闭。。。")
}

// 订阅模式 Step1：创建 RabbitMQ 实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 和普通模式创建不同，这里不指定队列名，而是指定交换机名
	rabbitMQ := NewRabbitMq("", exchangeName, "")
	return rabbitMQ
}

// 订阅模式 Step2：生产消息
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		/*
			交换机类型：
				direct Exchange：将消息中的Routing key与该Exchange关联的所有Binding中的Routing key进行比较，如果相等，则发送到该Binding对应的Queue中。

				topic Exchange：将消息中的Routing key与该Exchange关联的所有Binding中的Routing key进行对比，如果匹配上了，则发送到该Binding对应的Queue中。

				fanout Exchange：直接将消息转发到所有binding的对应queue中，这种exchange在路由转发的时候，忽略Routing key。

				headers Exchange：将消息中的headers与该Exchange相关联的所有Binging中的参数进行匹配，如果匹配上了，则发送到该Binding对应的Queue中。
		*/
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "声明创建交换机失败")

	// 2.发送消息
	e = r.channel.Publish(
		r.Exchange,
		// 在 pub/sub 订阅模式下，这里的key要为空
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

}

// 订阅模式 Step3:消费端代码
func (r *RabbitMQ) ReceiverSub() {
	// 1、和生产一样，首先尝试创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		//  交换机类型
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "消费端尝试创建交换机失败")

	// 2、尝试创建队列，队列名为空，随机生成队列名
	queue, e := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(e, "消费端尝试创建随机名队列失败")

	// 3、将创建的队列绑定到交换机 exchange 中
	e = r.channel.QueueBind(
		queue.Name,
		// 在 pub/sub 订阅模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)
	r.failOnError(e, "消费端绑定队列和交换机失败")

	messages, e := r.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	DonePubSub = make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("接受到订阅模式下的消息：%s", d.Body)
		}
	}()

	log.Printf("订阅模式 Pub/Sub 消费者已开启,队列名:%s，等待消息产生。。。", queue.Name)
	<-DonePubSub

	r.Destory()
	log.Printf("订阅模式消费者关闭。。。")
}

// 路由模式 Step1:创建实例
func NewRabbitMQRouting(exchangeName, key string) *RabbitMQ {
	//路由模式，相对于订阅模式，需要交换机，并且还有路由 key
	return NewRabbitMq("", exchangeName, key)
}

// 路由模式 Step2:生产消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.老样子，先尝试创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		//
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "路由模式创建交换机失败~")

	// 2.发送消息
	e = r.channel.Publish(
		r.Exchange,
		// 路由模式需要设置 key
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 路由模式 Step3:消费消息
func (r *RabbitMQ) ReceiverRouting() {
	// 1.试探性创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		//
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "路由模式接受消息创建交换机失败~")

	// 2.试探性创建队列，不填入队列名，随机生成
	queue, e := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(e, "路由模式接受消息创建队列失败")

	// 3.将队列绑定到交换机中，并且需要加入路由 key
	e = r.channel.QueueBind(
		queue.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	msgs, e := r.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	DoneRouting = make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("路由模式下接受到的消息：%s", d.Body)
		}
	}()

	log.Printf("路由模式 routing 消费者已开启,队列名:%s，等待消息产生。。。", queue.Name)
	<-DoneRouting

	r.Destory()
	log.Printf("路由模式消费者关闭。。。")
}

// Topic模式 Step1:创建实例
func NewRabbitMQTopic(exchangeName, key string) *RabbitMQ {
	//Topic模式，创建和订阅模式一样，需要交换机，并且还有路由 key
	return NewRabbitMq("", exchangeName, key)
}

// Topic模式 Step2:生产消息
func (r *RabbitMQ) PublishTopic(message string) {
	// 1.老样子，先尝试创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		//这里的类型要换成 topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "Topic模式创建交换机失败~")

	// 2.发送消息
	e = r.channel.Publish(
		r.Exchange,
		// 路由模式需要设置 key
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// Topic模式 Step3:消费消息
/*
	routing_key:必须是由点隔开的一系列的标识符组成。标识符可以是任何东西，但是一般都与消息的某些特性相关

		*可以匹配一个标识符。
		#可以匹配0个或多个标识符
*/
func (r *RabbitMQ) ReceiverTopic() {
	// 1.试探性创建交换机
	e := r.channel.ExchangeDeclare(
		r.Exchange,
		//
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(e, "Topic模式接受消息创建交换机失败~")

	// 2.试探性创建队列，不填入队列名，随机生成
	queue, e := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(e, "Topic模式式接受消息创建队列失败")

	// 3.将队列绑定到交换机中，并且需要加入路由 key
	e = r.channel.QueueBind(
		queue.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	msgs, e := r.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	DoneRouting = make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Topic模式下接受到的消息：%s", d.Body)
		}
	}()

	log.Printf("Topic模式 消费者已开启,队列名:%s，等待消息产生。。。", queue.Name)
	<-DoneRouting

	r.Destory()
	log.Printf("Topic模式消费者关闭。。。")
}
