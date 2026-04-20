package mqtt

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"sensor/internal/config"
	"sensor/internal/service"
)

// MQTTConsumer 表示MQTT消费者
type MQTTConsumer struct {
	client    mqtt.Client            // MQTT客户端
	config    config.MQTTConfig      // MQTT配置
	service   *service.SensorService // 传感器服务
	connected bool                   // 连接状态
}

// MQTTConfig 是config.MQTTConfig的别名
type MQTTConfig = config.MQTTConfig

// NewConsumer 创建MQTT消费者
func NewConsumer(cfg config.MQTTConfig, svc *service.SensorService) *MQTTConsumer {
	return &MQTTConsumer{
		config:  cfg,
		service: svc,
	}
}

// Start 启动MQTT消费者
func (c *MQTTConsumer) Start() error {
	opts := mqtt.NewClientOptions().
		AddBroker(c.config.Broker).
		SetClientID(c.config.ClientID).
		SetCleanSession(true).
		SetConnectionLostHandler(c.onConnectionLost).
		SetOnConnectHandler(c.onConnect).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
	}
	if c.config.Password != "" {
		opts.SetPassword(c.config.Password)
	}

	c.client = mqtt.NewClient(opts)

	// 连接到broker
	token := c.client.Connect()
	if token.Error() != nil {
		log.Printf("MQTT Broker 连接失败")
		c.connected = false
		return token.Error()
	}

	log.Printf("MQTT消费者已启动，连接到 %s", c.config.Broker)
	return nil
}

// Stop 停止MQTT消费者
func (c *MQTTConsumer) Stop() {
	if c.client != nil && c.connected {
		c.client.Disconnect(5000)
		log.Println("MQTT消费者已停止")
	}
}

// onConnectionLost 当连接丢失时调用
func (c *MQTTConsumer) onConnectionLost(client mqtt.Client, err error) {
	log.Printf("MQTT连接丢失: %v", err)
	c.connected = false
}

// onConnect 当连接建立时调用
func (c *MQTTConsumer) onConnect(client mqtt.Client) {
	log.Println("MQTT已连接")

	// 订阅主题
	topic := c.config.Topic
	qos := byte(c.config.QOS)

	token := client.Subscribe(topic, qos, c.messageHandler)
	if token.Wait() && token.Error() != nil {
		log.Printf("订阅主题 %s 失败: %v", topic, token.Error())
		return
	}

	c.connected = true
	log.Printf("已订阅主题: %s", topic)
}

// messageHandler 处理接收到的MQTT消息
func (c *MQTTConsumer) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("收到主题 %s 的消息", msg.Topic())

	if err := c.service.ProcessMQTTMessage(msg.Payload()); err != nil {
		log.Printf("处理消息失败: %v", err)
	}
}

// NewMQTTConfig 从config.MQTTConfig创建MQTTConfig
func NewMQTTConfig(cfg config.MQTTConfig) config.MQTTConfig {
	return cfg
}

// StartWithGracefulShutdown 启动支持优雅关闭的MQTT消费者
func StartWithGracefulShutdown(cfg config.MQTTConfig, svc *service.SensorService) error {
	consumer := NewConsumer(cfg, svc)

	if err := consumer.Start(); err != nil {
		return err
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("正在关闭MQTT消费者...")
	consumer.Stop()

	return nil
}
