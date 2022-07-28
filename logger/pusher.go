package logger

import (
	"log_consumer/rabbit"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/streadway/amqp"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var once sync.Once

const InfoLogLevel = "info"
const WarnLogLevel = "warn"
const ErrorLogLevel = "error"
const FatalLogLevel = "critical"
const DebugLogLevel = "debug"

// Pusher defined type with interface
type Pusher interface {
	Info(body map[string]interface{}) error
	Warn(body map[string]interface{}) error
	Error(body map[string]interface{}) error
	Fatal(body map[string]interface{}) error
	Debug(body map[string]interface{}) error

	push(queue string, body map[string]interface{}) error

	InitPusher(client rabbit.Client)
	DefinePid(host string, action string)
}

type pusher struct {
	Client rabbit.Client
	UniquePid string
	Action string
	sync.RWMutex
}

var instance *pusher = nil

func GetInstance() Pusher {
	once.Do(func() {
		instance = new(pusher)
	})
	return instance
}

func (p *pusher) DefinePid(host string, action string) {
	p.Lock()
	defer p.Unlock()
	randomObj := rand.New(rand.NewSource(time.Now().UnixNano()))
	pid := strconv.Itoa(randomObj.Intn(100)) + host
	p.UniquePid = b64.StdEncoding.EncodeToString([]byte(pid))
	p.Action = action
}

func (p *pusher) InitPusher(client rabbit.Client) {
	p.Lock()
	defer p.Unlock()
	p.Client = client
}

func (p *pusher) Info(body map[string]interface{}) error {
	p.RLock()
	defer p.RUnlock()
	return p.push(rabbit.InfoQueue, body)
}

func (p *pusher) Error(body map[string]interface{}) error {
	p.RLock()
	defer p.RUnlock()
	return p.push(rabbit.ErrorQueue, body)
}

func (p *pusher) Warn(body map[string]interface{}) error {
	p.RLock()
	defer p.RUnlock()
	return p.push(rabbit.WarnQueue, body)
}

func (p *pusher) Fatal(body map[string]interface{}) error {
	p.RLock()
	defer p.RUnlock()
	return p.push(rabbit.FatalQueue, body)
}

func (p *pusher) Debug(body map[string]interface{}) error {
	p.RLock()
	defer p.RUnlock()
	return p.push(rabbit.DebugQueue, body)
}

func (p *pusher) push(queue string, body map[string]interface{}) error  {
	jsonString, parseErr := json.Marshal(body)
	if parseErr != nil {
		return parseErr
	}
	msgBody := rabbit.MessageBody{
		Body: jsonString,
		Pid: p.UniquePid,
		CreatedTime: time.Now().UnixNano(),
		Action: p.Action,
		Level: defineLevel(queue),
	}
	fullBody, err := msgBody.Encode()
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        fullBody,
	}
	return p.Client.Chanel.Publish("", queue, false, false, msg)
}

func defineLevel(queue string) string {
	switch queue {
	case rabbit.InfoQueue:
		return InfoLogLevel
	case rabbit.WarnQueue:
		return WarnLogLevel
	case rabbit.ErrorQueue:
		return ErrorLogLevel
	case rabbit.FatalQueue:
		return FatalLogLevel
	case rabbit.DebugQueue:
		return DebugLogLevel
	default:
		return ""
	}
}
