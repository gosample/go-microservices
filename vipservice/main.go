package main

import (
	"flag"
	"fmt"
	"github.com/manachyn/go-microservices-blog/vipservice/messaging"
	"github.com/manachyn/go-microservices-blog/vipservice/service"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
        "github.com/manachyn/go-microservices-blog/vipservice/config"
        "os"
        "os/signal"
        "syscall"
)

var appName = "vipservice"

var consumer messaging.IMessagingConsumer

func main() {
	fmt.Println("Starting " + appName + "...")
	parseFlags()

        config.LoadConfiguration(viper.GetString("configServerUrl"), appName, viper.GetString("profile"))
        initializeMessaging()

        // Call the subscribe method with queue name and callback function
	consumer.Subscribe("vipQueue", onMessage)

        // Makes sure connection is closed when service exits.
        handleSigterm(func() {
                if consumer != nil {
                        consumer.Close()
                }
        })

        service.StartWebServer(viper.GetString("server_port"))
}

func onMessage(delivery amqp.Delivery) {
	fmt.Printf("Got a message: %v\n", string(delivery.Body))
}

func parseFlags() {
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerUrl := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")

	flag.Parse()
	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
}

func initializeMessaging() {
	if !viper.IsSet("broker_url") {
		panic("No 'broker_url' set in configuration, cannot start")
	}
	consumer = &messaging.MessagingConsumer{}
	consumer.ConnectToBroker(viper.GetString("broker_url"))
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        signal.Notify(c, syscall.SIGTERM)
        go func() {
                <-c
                handleExit()
                os.Exit(1)
        }()
}
