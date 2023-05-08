package infrastructure

import (
	"encoding/json"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/application"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/model"
)

func ConsumeMessages(app *application.App, queueName string) {
	dep := application.SetupDependencyInjection(app)

	// Subscribing to queues for getting messages.
	messages, err := app.RabbitMQ.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		app.Logger.Error("Failed consume message in queue", queueName)
	}

	app.Logger.Info("Waiting Message in Queues ", queueName, ".....")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for msg := range messages {
			switch queueName {
			case app.Config.RabbitMQ.QueueCreateShortener:
				var req model.CreateShortRequest

				err := json.Unmarshal(msg.Body, &req)
				if err != nil {
					app.Logger.Error("Unmarshal JSON ERROR, ", err)
				}

				app.Logger.Info("Success Consume Message :", req)

				err = dep.ShortController.ProcessCreateShortUser(app.Context, &req)
				if err != nil {
					app.Logger.Error("ProcessCreateShortUser ERROR, ", err)
				}

				app.Logger.Info("Success Process Message : ", req)
			case app.Config.RabbitMQ.QueueUpdateVisitor:
			}
		}
	}()

	<-forever
}
