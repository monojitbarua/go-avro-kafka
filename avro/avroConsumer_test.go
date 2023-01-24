package kafka

import (
	"testing"
)

const GROUP_ID string = "go-kafka-avro-consumer-client"

func TestAvroConsumer_ProcessAvroMsg(t *testing.T) {

	callbacks := &ConsumerCallbacks{}

	avroConsumer, err := NewAvroConsumer(brokers, schemas, "EMPLOYEE_AVRO_TOPIC", GROUP_ID, *callbacks)
	if err != nil {
		panic(1)
	}

	avroConsumer.Consume()
}
