package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	// 生产者往kafka发送数据的模式(3种) 1.吧数据发送给leader就成功， 效率最高，安全性最低 2：吧数据发送给leader，等待leader回ack 3：吧数据发给leader 确保follower从leader拉去数据恢复ack给leader， leader再恢复ack；安全性最高
	config.Producer.RequiredAcks = sarama.WaitForAll
	// kafka选择分区的模式（3种） 1 指定往那个分区写2 指定key kafka根据key做hash然后决定那个分区 3 轮询
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	// 构造一个消息体
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is atest log")
	// 链接kafka
	client, err := sarama.NewSyncProducer([]string{"192.168.1.7:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err: ", err)
		return
	}
	fmt.Println("pid:%v offset:%v\n", pid, offset)
}
