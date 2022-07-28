package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log_consumer/log_handler"
	"log_consumer/mongo"
	"log_consumer/rabbit"
	"strconv"
	"sync"
)

const queueCount = 5

type Configurations struct {
	Config Config `yaml:"config"`
	Mongo map[string]interface{} `yaml:"mongo"`
}

type Config struct {
	RabbitHostPort string `yaml:"rabbit_host_port"`
	RabbitUser string `yaml:"rabbit_user"`
	RabbitPassword string `yaml:"rabbit_pass"`
	MongoHostPort string `yaml:"mongo_host_port"`
	MongoUser string `yaml:"mongo_user"`
	MongoPassword string `yaml:"mongo_pass"`
}

type ConsumersInfo struct {
	QueueName string
	ConsumersCount int
}

type Consumer struct {
	Client *rabbit.Client
	Queue  amqp.Queue
}

var RClient rabbit.Client

var consumersArray [queueCount]ConsumersInfo
var mainError error

var LogConsumer []Consumer

func init() {
	yfile, parseErr := ioutil.ReadFile("configurations.yml")
	if parseErr != nil {panicIfNeed(parseErr, "")}
	configurations := Configurations{}
	parseErr = yaml.Unmarshal(yfile, &configurations)
	if parseErr != nil {panicIfNeed(parseErr, "")}

	mainError = RClient.SetConnection(
		configurations.Config.RabbitUser,
		configurations.Config.RabbitPassword,
		configurations.Config.RabbitHostPort,
	)
	panicIfNeed(mainError, "connect to rabbit")

	mainError = mongo.Connect(
		configurations.Config.MongoHostPort,
		configurations.Config.MongoUser,
		configurations.Config.MongoPassword,
		configurations.Mongo,
	)
	panicIfNeed(mainError, "connect to mongo")

	consumersArray = [queueCount]ConsumersInfo{
		{QueueName: rabbit.InfoQueue, ConsumersCount: 4},
		{QueueName: rabbit.ErrorQueue, ConsumersCount: 2},
		{QueueName: rabbit.WarnQueue, ConsumersCount: 1},
		{QueueName: rabbit.FatalQueue, ConsumersCount: 1},
		{QueueName: rabbit.DebugQueue, ConsumersCount: 1},
	}

	LogConsumer = make([]Consumer, queueCount)
}

// consumer
func main() {
	var wg sync.WaitGroup
	// declare all queues
	for i, consumerInfo := range consumersArray {
		LogConsumer[i] = Consumer{}
		LogConsumer[i].Client = &RClient
		mainError = LogConsumer[i].CreateQueue(consumerInfo.QueueName)
		panicIfNeed(mainError, "declare queue: " + consumerInfo.QueueName)
		for j := 0; j < consumerInfo.ConsumersCount; j++ {
			wg.Add(1)
			go LogConsumer[i].Execute(&wg, consumerInfo.QueueName + "_" + strconv.Itoa(j))
		}
	}
	wg.Wait()

	mainError = RClient.CloseConnection()
	panicIfNeed(mainError, "log_consumer stopped")
}

func (consumer *Consumer) Execute(wg *sync.WaitGroup, name string) {
	defer wg.Done()
	var distributeErr error
	msgs, err := consumer.Client.Chanel.Consume(
		consumer.Queue.Name, // queue
		name,                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		true,               // no-wait
		nil,                 // args
	)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			distributeErr = log_handler.Distribute(msg.Body)
			if distributeErr != nil {
				continue
			}
		}
	}()
	<-forever
}

func (consumer *Consumer) CreateQueue(name string) error {
	var err error
	consumer.Queue, err = consumer.Client.Chanel.QueueDeclare(
		name, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	return err
}

func panicIfNeed(err error, successMessage string) {
	if err != nil {
		panic(err)
	} else {
		fmt.Println(successMessage)
	}
}
