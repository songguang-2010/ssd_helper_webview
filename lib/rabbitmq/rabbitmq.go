package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//消息结构
type MSG struct {
	Body    []byte
	Tag     uint64
	Channel string
}

//公共配置参数结构
type MQConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
	Exchange string
	Key      string
}

//对象共有属性
type RabbitmqHandler struct {
	config     *MQConf
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	key        string
}

//公共对象方法
func (h *RabbitmqHandler) SetConfig(mqConfig *MQConf) error {
	h.config = mqConfig
	return nil
}

func (h *RabbitmqHandler) Connect() error {
	msn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", h.config.User, h.config.Password, h.config.Host, h.config.Port, h.config.Vhost)

	connect, err := amqp.Dial(msn)
	if err != nil {
		return err
	}
	h.connection = connect

	channel, err := h.connection.Channel()
	if err != nil {
		return err
	}
	h.channel = channel

	if err = h.channel.Qos(1, 0, false); err != nil {
		return nil
	}

	// var args amqp.Table

	err = h.channel.ExchangeDeclare(h.config.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}
	h.exchange = h.config.Exchange

	return nil
}

func (h *RabbitmqHandler) Clean() error {
	if h.channel != nil {
		if err := h.channel.Close(); err != nil {
			return err
		}
	}

	if h.connection != nil {
		if err := h.connection.Close(); err != nil {
			return err
		}
	}

	return nil
}
