package rabbitmq

import (
	"github.com/streadway/amqp"
)

//发布对象配置参数结构
type MQPubConf struct {
}

//发布者对象
type RabbitmqPubHandler struct {
	RabbitmqHandler
	configPub *MQPubConf
}

//发布对象方法
func (ph *RabbitmqPubHandler) SetConfig(mqConfig *MQConf, mqPubConfig *MQPubConf) error {
	err := ph.RabbitmqHandler.SetConfig(mqConfig)
	if err != nil {
		return err
	}
	ph.configPub = mqPubConfig
	return nil
}

func (ph *RabbitmqPubHandler) Connect() error {
	err := ph.RabbitmqHandler.Connect()
	if err != nil {
		return err
	}
	ph.key = ph.config.Key
	return nil
}

// 发布者方法
func (ph *RabbitmqPubHandler) Publish(msg []byte) (*RabbitmqPubHandler, error) {
	if err := ph.channel.Publish(ph.config.Exchange, ph.config.Key, false, false, amqp.Publishing{ContentType: "text/plain", Body: msg}); err != nil {
		return ph, err
	}
	return ph, nil
}

func CreatePublishHandler(mqConfig *MQConf, mqPubConfig *MQPubConf) (*RabbitmqPubHandler, error) {
	pubHandler := &RabbitmqPubHandler{}
	err := pubHandler.SetConfig(mqConfig, mqPubConfig)
	if err != nil {
		return pubHandler, err
	}
	err = pubHandler.Connect()
	if err != nil {
		return pubHandler, err
	}

	return pubHandler, nil
}
