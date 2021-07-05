package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "time"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}

func main() {
	bot, err := linebot.New(
		os.Getenv("LINEBOT_CHANNEL_SECRET"),
		os.Getenv("LINEBOT_CHANNEL_TOKEN"),
		// "a73b62d06a29b77d3b57b3d3b0aa0e7b",
		// "VhC7qpsC9Op/QN1MDc61EGAN5Jqiq2fl5RlyzGZjVJr0CnZE7gs2G52HOt9pWPEzFYvY74eRqzC939lWERLSxYZk1uaFMSQpy0v92hjZfVvyFoOX9VzMSAULznGrP5sa5wE+viP8gkG2d939jxiV3QdB04t89/1O/w1cDnyilFU=",
	  )
	  if err != nil {
		log.Fatal(err)
	  }
	  router := gin.Default()

	
	  router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
		  if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		  } else {
			c.Writer.WriteHeader(500)
		  }
		  return
		}
		for _, event := range events {
		  if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			  case *linebot.TextMessage:
				if message.Text=="ข้อความ"{
				  if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				  }
				}else if message.Text=="ปิดไฟ"{
					mqtt_main("relay1_off");
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ปิดไฟเรียบร้อย")).Do(); err != nil {
						log.Print(err)
					}
				}else if message.Text=="เปิดไฟ"{
					mqtt_main("relay1_on");
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("เปิดไฟเรียบร้อย")).Do(); err != nil {
						log.Print(err)
					}
				}
			  
			}
		  }
		}
	  })
	  router.Run(":" + os.Getenv("PORT"))
	//   router.Run(":5600")
	
}


func mqtt_main(msg_line string){
	var broker = "soldier.cloudmqtt.com"
    var port = 10174
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("brdhfcif")
    opts.SetPassword("gviTCGqRHgB9")
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    sub(client)
    publish(client,msg_line)

    client.Disconnect(250)
}

func publish(client mqtt.Client,msg string) {
        token := client.Publish("/ESP/LED", 0, false, msg)
        token.Wait()
        time.Sleep(time.Second)
}

func sub(client mqtt.Client) {
    topic := "/ESP/LED"
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
  fmt.Printf("Subscribed to topic: %s", topic)
}


//https://www.emqx.io/blog/how-to-use-mqtt-in-golang
