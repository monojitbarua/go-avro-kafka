package kafka

import (
	"fmt"
	"testing"
	"time"
)

const BROKER_IP string = "localhost:9092"
const SCHEMA_REGISTRY_IP string = "http://localhost:8081"

var brokers = []string{BROKER_IP}
var schemas = []string{SCHEMA_REGISTRY_IP}

func Test_GoTimeout(t *testing.T) {

	timeoutChannel := make(chan bool, 1)
	go func() {
		time.Sleep(30 * time.Second)
		timeoutChannel <- false
	}()

	go func() {
		time.Sleep(300 * time.Second)
		timeoutChannel <- false
	}()

	if <-timeoutChannel {
		return
	}

	// timeout := time.After(300 * time.Second)
	// done := make(chan bool)
	// go func() {
	// 	time.Sleep(35 * time.Second)
	// 	done <- true
	// }()

	// select {
	// case <-timeout:
	// 	t.Fatal("Test didn't finish in time")
	// case <-done:
	// }
}

func TestAvroProducer_Add(t *testing.T) {

	t.Setenv("timeout", "30m")

	schemaRegistryClient := NewCachedSchemaRegistryClient(schemas)
	codec, err := schemaRegistryClient.GetSchema(1)
	if err != nil {
		panic(1)
	}

	avroProducer, err := NewAvroProducer(brokers, schemas)
	if err != nil {
		panic(2)
	}

	for i := 1; i <= 10; i++ {
		time.Sleep(10 * time.Second)
		v := fmt.Sprintf(`{"val":%d}`, i)
		partition, offset, err := avroProducer.Add("EMPLOYEE_AVRO_TOPIC", codec.Schema(), []byte("key"), []byte(v))
		fmt.Println("-------- Message has been published to partition: ", partition, " offset:", offset, " error:", err, " schema: ", codec.Schema())
	}

}
