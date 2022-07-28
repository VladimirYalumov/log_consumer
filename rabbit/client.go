package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

const InfoQueue = "deewave_info"
const WarnQueue = "deewave_warn"
const ErrorQueue = "deewave_error"
const FatalQueue = "deewave_fatal"
const DebugQueue = "deewave_debug"

type Client struct {
	Connection *amqp.Connection
	Chanel *amqp.Channel
}

type MessageBody struct {
	CreatedTime  int64 `json:"created_time"`
	Pid     string `json:"pid"`
	Body    []byte `json:"body"`
	Action  string `json:"action"`
	Level string `json:"level"`
}

func (s *Client) SetConnection(user string, password string, host string) error {
	var err error
	s.Connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, password, host))
	if err != nil {
		return err
	}
	s.Chanel, err = s.Connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (s *Client) CloseConnection() error {
	err := s.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (msg *MessageBody) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *MessageBody) Decode(encodeString []byte) error {
	return json.Unmarshal(encodeString, msg)
}
