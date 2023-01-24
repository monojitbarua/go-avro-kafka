package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	kafka "github.com/monojitbarua/go-avro-kafka/avro"
	"github.com/monojitbarua/go-avro-kafka/util"
)

var conf map[string]string = util.Conf
var brokers = []string{conf["bootstrap.servers"]}
var schemas = []string{conf["schema.registry"]}

func main() {

	//start consumer
	go func() { consume() }()

	//start a rest end point which will receive audit log and publish
	r := mux.NewRouter()
	r.HandleFunc("/produce/message", POST_HANDLER)
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		fmt.Println("Starting Server ...")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	// Graceful Shutdown
	waitForShutdown(server)
}

func POST_HANDLER(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	partition, offset := publish(string(reqBody))
	json.NewEncoder(w).Encode(fmt.Sprintf("partition:%d, offset:%d", partition, offset))
}

func waitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	fmt.Println("Shutting down")
	os.Exit(0)
}

func publish(auditlog string) (partition int32, offset int64) {
	schemaRegistryClient := kafka.NewCachedSchemaRegistryClient(schemas)

	schemaId, err := strconv.Atoi(conf["avro.schema.id"])
	if err != nil {
		fmt.Printf("\nError in getting schema id: %s", err)
		panic(0)
	}

	codec, err := schemaRegistryClient.GetSchema(schemaId)
	if err != nil {
		fmt.Printf("\nError in schema registry connection: %s", err)
		panic(1)
	}

	avroProducer, err := kafka.NewAvroProducer(brokers, schemas)
	if err != nil {
		fmt.Printf("\nError in avro producer: %s", err)
		panic(2)
	}

	partition, offset, err = avroProducer.Add(conf["topic.name"], codec.Schema(), []byte("key"), []byte(auditlog))
	if err != nil {
		fmt.Printf("\nError in avro publish: %s", err)
		panic(3)
	}
	fmt.Printf("\nMessage has been published to topic: %s, partition: %d, offset: %d", conf["topic.name"], partition, offset)
	return partition, offset
}

func consume() {

	callbacks := &kafka.ConsumerCallbacks{}

	avroConsumer, err := kafka.NewAvroConsumer(brokers, schemas, conf["topic.name"], conf["group.id"], *callbacks)
	if err != nil {
		fmt.Printf("\nError in consumer: %s broker: %s", err, conf["bootstrap.servers"])
		panic(1)
	}
	avroConsumer.Consume()
}
