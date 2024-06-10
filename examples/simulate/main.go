package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type message struct {
	body  string
	line  uint
	level string
}

var (
	faker  = gofakeit.NewCrypto()
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type logEvent struct {
	SourceHost  string    `json:"source_host"`
	Method      string    `json:"method"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"@timestamp"`
	File        string    `json:"file"`
	Application string    `json:"application"`
	UserId      int       `json:"user_id"`
	LineNumber  string    `json:"line_number"`
	ThreadName  string    `json:"thread_name"`
	Version     int       `json:"@version"`
	LoggerName  string    `json:"logger_name"`
	Class       string    `json:"class"`
}

func threadName() string {
	f := rand.Float64()
	switch {
	case f < 0.01:
		return "main"
	case f < 0.30:
		return fmt.Sprintf("catalina-worker-%d", random.Intn(100)+1)
	case f < 0.70:
		return fmt.Sprintf("pool-%d-thread-%d", random.Intn(2)+1, random.Intn(10)+1)
	case f < 0.99:
		return fmt.Sprintf("http-nio-%d-exec-%d", random.Intn(10)+1, random.Intn(10)+1)
	default:
		return fmt.Sprintf("%s-%d", faker.VerbAction(), random.Intn(10)+1)
	}
}

var classNames = map[string][]string{
	"greedy-turnip": {
		"com.gaggle.greedy.turnip.Model",
		"com.gaggle.greedy.turnip.Controller",
		"com.gaggle.greedy.turnip.Service",
		"com.gaggle.greedy.turnip.Repository",
	},
	"mangled-horseshoe": {
		"com.gaggle.mangled.horseshoe.service.Impl",
		"com.gaggle.mangled.horseshoe.impl.Service",
		"com.gaggle.mangled.horseshoe.Action",
	},
	"funky-cherub": {
		"com.gaggle.funky.cherub.App",
		"com.gaggle.funky.cherub.Logger",
		"com.gaggle.funky.cherub.Serializer",
		"com.gaggle.funky.cherub.util.Util",
	},
}
var methodNames = map[string][]string{
	"com.gaggle.greedy.turnip.Model": {
		"setAge",
	},
	"com.gaggle.greedy.turnip.Controller": {
		"initBook",
		"getBook",
		"updateBook",
		"deleteBook",
		"getBooks",
	},
	"com.gaggle.greedy.turnip.Service": {
		"persistBook",
		"removeBook",
		"validateBook",
	},
	"com.gaggle.greedy.turnip.Repository": {
		"recoverBook",
		"restoreTurnip",
		"findBook",
		"findTurnip",
	},
	"com.gaggle.mangled.horseshoe.service.Impl": {
		"addShoe",
		"removeShoe",
		"findHorse",
		"findShoe",
	},
	"com.gaggle.mangled.horseshoe.impl.Service": {
		"getAddress",
		"setAddress",
		"getHorse",
		"setHorse",
	},
	"com.gaggle.mangled.horseshoe.Action": {
		"mangleShoe",
		"unmangleShoe",
		"removeNail",
		"addNail",
	},
	"com.gaggle.funky.cherub.App": {
		"start",
		"stop",
		"restart",
		"pause",
		"resume",
	},
	"com.gaggle.funky.cherub.Logger": {
		"invertCherub",
		"convertCherub",
		"revertCherub",
		"assertCherub",
	},
	"com.gaggle.funky.cherub.util.Util": {
		"cherubimAsSerafim",
		"serafimAsCherubim",
		"firstCherub",
		"lastCherub",
	},
}

var messages = map[string]map[string][]message{}

func init() {
	staticFaker := gofakeit.New(101)
	for class, methods := range methodNames {
		for _, method := range methods {
			for i := 0; i < 10; i++ {
				if messages[class] == nil {
					messages[class] = map[string][]message{}
				}

				var level string
				r := random.Float64()
				switch {
				case r < 0.0001:
					level = "FATAL"
				case r < 0.001:
					level = "ERROR"
				case r < 0.01:
					level = "WARN"
				case r < 0.7:
					level = "INFO"
				case r < 0.999:
					level = "DEBUG"
				default:
					level = "TRACE"
				}

				sentence := staticFaker.Sentence(10)
				if class == "com.gaggle.funky.cherub.Serializer" {
					if ba, err := faker.JSON(
						&gofakeit.JSONOptions{
							Type:     "object",
							RowCount: 1,
							Fields: []gofakeit.Field{
								{
									Name:     "widgetCount",
									Function: "number",
								},
								{
									Name:     "action",
									Function: "verb",
								},
								{
									Name:     "model",
									Function: "noun",
								},
							},
						},
					); err == nil {
						sentence = string(ba)
					}
				}
				messages[class][method] = append(
					messages[class][method], message{
						body:  sentence,
						line:  uint(staticFaker.Uint64()%10_000 + 1),
						level: level,
					},
				)
			}
		}
	}
}

func generateLogEvent() logEvent {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	application := []string{"greedy-turnip", "mangled-horseshoe", "funky-cherub"}[random.Intn(3)]
	class := classNames[application][random.Intn(len(classNames[application]))]
	method := methodNames[class][random.Intn(len(methodNames[class]))]
	msg := messages[class][method][random.Intn(len(messages[class][method]))]
	file := path.Base(strings.ReplaceAll(class, ".", "/") + ".java")

	return logEvent{
		SourceHost:  faker.IPv4Address(),
		Method:      method,
		Level:       msg.level,
		Message:     msg.body,
		Environment: []string{"dev", "staging", "prod"}[random.Intn(3)],
		Timestamp:   time.Now(),
		Application: application,
		UserId:      0,
		LineNumber:  strconv.Itoa(int(msg.line)),
		ThreadName:  threadName(),
		Version:     1,
		File:        file,
		Class:       class,
		LoggerName:  class,
	}
}

func sendEvent(svc *lambda.Client, n int) error {
	payload, err := json.Marshal(generateLogEvent())
	if err != nil {
		return err
	}

	// Invoke the Lambda function
	input := &lambda.InvokeInput{
		FunctionName: aws.String("dev-spa7aa-couture-simulator-Simulate-" + strconv.Itoa(n)),
		Payload:      payload,
	}

	result, err := svc.Invoke(context.TODO(), input)
	if err != nil {
		return err
	}

	fmt.Println(string(result.Payload))
	return nil
}

func main() {
	const workerCount = 9

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		panic(err)
	}
	svc := lambda.NewFromConfig(cfg)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	running := true
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	for w := 0; w < workerCount; w++ {
		go func() {
			defer wg.Done()
			id := w + 1
			for running {
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
				err := sendEvent(svc, id)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	go func() {
		<-sigs
		running = false
		wg.Wait()
	}()

	wg.Wait()
}
