package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptrace"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"

	"blue-admin.com/configs"
	"blue-admin.com/observe"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
)

type sample_message struct{}

func RabbitConsumer(queue_name string, env string) {

	// Loading configuration file
	configs.AppConfig.SetEnv(env)

	//  tracer
	// tp := observe.InitTracer()
	// defer func() {
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Printf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	// Getting app connection and channel
	connection, channel, err := QeueConnect(queue_name)
	if err != nil {
		fmt.Println("Failed to establish connection:", err)
		return
	}
	defer connection.Close()
	defer channel.Close()

	// ########################################
	// Declaring consumer with its properties over the channel opened
	msgs, err := channel.Consume(
		queue_name, // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)

	// ###########################################

	if err != nil {
		fmt.Println("Failed to consume messages:", err)
		return
	}

	// Process received messages based on their types
	// Using a goroutine for asynchronous message consumption
	go func(msg <-chan amqp.Delivery) {
		ctx := context.Background()

		for msg := range msgs {
			// Extract the span context out of the AMQP header.

			switch msg.Type {
			case "BULK_MAIL": // make sure provide the type in the published message so to switch
				var message sample_message
				err := json.Unmarshal(msg.Body, &message)
				if err != nil {
					fmt.Println("Failed to unmarshal message:", err)
					msg.Reject(true)
					break
				}
				msg.Ack(true)
			case "REQUEST":
				//  Parsing Request object
				var reqData RequestObject
				err := json.Unmarshal(msg.Body, &reqData)
				if err != nil {
					fmt.Println(err.Error())
				}

				// extracting request from otel propagator
				propagator := propagation.TraceContext{}
				ctx = propagator.Extract(ctx, propagation.HeaderCarrier(reqData.Tp))

				//starting span for http client request
				tracer, span := observe.AppTracer.Start(ctx, fmt.Sprintf("started-esb%v", rand.Intn(1000)))

				// client with otel http middleware
				client := http.Client{
					Transport: otelhttp.NewTransport(
						http.DefaultTransport,
						// By setting the otelhttptrace client in this transport, it can be
						// injected into the context after the span is started, which makes the
						// httptrace spans children of the transport one.
						otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
							return otelhttptrace.NewClientTrace(ctx)
						}),
					),
				}

				// Build request to be sent
				req, rerr := http.NewRequestWithContext(tracer, reqData.Method, fmt.Sprintf("http://%v%v", reqData.Host, reqData.Endpoint), nil)
				if rerr != nil {
					fmt.Printf("failed to perform request ######: %v\n", rerr)
				}

				// generating uuid
				gen, _ := uuid.NewV7()
				id := gen.String()

				//  performing the request
				req_body := fmt.Sprintf("Method: %v\t %v \nBody: %v\n", req.Method, req.URL, req.Body)
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("failed to perform request: %v\n", err)
					span.SetAttributes(attribute.String("esb-id", id))
					span.SetAttributes(attribute.String("esb-request", req_body))
					span.SetAttributes(attribute.String("esb-error", err.Error()))
					msg.Reject(true)
					// msg.Ack(true)
					span.End()
					break
				}

				body, _ := io.ReadAll(resp.Body)
				span.SetAttributes(attribute.String("esb-id", id))
				span.SetAttributes(attribute.String("esb-request", req_body))
				span.SetAttributes(attribute.String("esb-response", string(body)))
				msg.Ack(true)
				span.End()

			default:
				fmt.Println("Unknown Task Type")
			}
		}
	}(msgs)

	fmt.Println("Waiting for messages...")
	select {}
}
