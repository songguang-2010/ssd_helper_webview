package rabbitmq

import (
	"github.com/streadway/amqp"
	"lib/coroutine"
	log "lib/logwrap"
)

//订阅对象配置参数结构
type MQSubConf struct {
	Queue string
}

//消费者对象
type RabbitmqSubHandler struct {
	RabbitmqHandler
	configSub *MQSubConf
	queue     string
}

//订阅对象方法
func (sh *RabbitmqSubHandler) SetConfig(mqConfig *MQConf, mqSubConfig *MQSubConf) error {
	err := sh.RabbitmqHandler.SetConfig(mqConfig)
	if err != nil {
		return err
	}
	sh.configSub = mqSubConfig
	return nil
}
func (sh *RabbitmqSubHandler) Connect() error {
	err := sh.RabbitmqHandler.Connect()
	if err != nil {
		return err
	}
	_, err = sh.channel.QueueDeclare(sh.configSub.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	sh.queue = sh.configSub.Queue

	if err = sh.channel.QueueBind(sh.configSub.Queue, sh.config.Key, sh.config.Exchange, false, nil); err != nil {
		return err
	}
	sh.key = sh.config.Key
	return nil
}

//订阅者方法
func (sh *RabbitmqSubHandler) Subscribe(callback func(MSG, *RabbitmqSubHandler)) (*RabbitmqSubHandler, error) {
	var msgs <-chan amqp.Delivery

	msgs, err := sh.channel.Consume(sh.configSub.Queue, "", false, false, false, false, nil)
	if err != nil {
		return sh, err
	}

	// go handleMsg(msgs, callback, sh)
	coroutine.Start(func(current int, total int) {
		handleMsg(msgs, callback, sh)
	}, 1)

	return sh, nil
}

//处理消息(顺序处理,如果需要多线程可以在回调函数中做手脚)
func handleMsg(msgs <-chan amqp.Delivery, callback func(MSG, *RabbitmqSubHandler), sh *RabbitmqSubHandler) {
	//catch panic error and free the resource
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Coroutine 'handleMsg' Runtime Error: ", err)
		}
	}()

	for d := range msgs {
		var msg MSG = MSG{
			Body: d.Body,
			Tag:  d.DeliveryTag,
		}
		callback(msg, sh)
	}
}

func (sh *RabbitmqSubHandler) Ack(m MSG, multiple bool) (err error) {
	sh.channel.Ack(m.Tag, multiple)
	return nil
}

func CreateSubscribeHandler(mqConfig *MQConf, mqSubConfig *MQSubConf) (*RabbitmqSubHandler, error) {
	subHandler := &RabbitmqSubHandler{}
	err := subHandler.SetConfig(mqConfig, mqSubConfig)
	if err != nil {
		return subHandler, err
	}
	err = subHandler.Connect()
	if err != nil {
		return subHandler, err
	}

	return subHandler, nil
}
