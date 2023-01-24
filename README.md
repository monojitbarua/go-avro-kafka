
# GO-AVRO-KAFKA

  

This is a code for kafka publisher and consumer using go-lang and avro as serializer and deserializer.

  
  

# Installation steps

In real life senario we will going to use the dev/stage/prod kafka instance. But here we will replicate those instace in local docker containers as below:

* Zookeeper

* Broker: kafka server

* Schema registry: stores schemas

* Control center: So that via web interface we can create/update/delete topic, schemas

- Start these instace after cloing the repo ([go-avro-kafka](https://github.com/monojitbarua/go-avro-kafka)) : `docker-compose up`

  
  

# Create schema using control center

* After successfully starting the above mentioned containers need to login in http://localhost:9021 and create topic `EMPLOYEE_AVRO_TOPIC`

* Then under schmea tab add below avro schema for audit log:

  

		{

		"type": "record",

		"name": "Employee",

		"fields": [{

		"name": "employeeName",

		"type": "string",

		"default": ""

		}, {

		"name": "department",

		"type": "string",

		"doc": ""

		}, {

		"name": "salary",

		"type": "float"

		},

		{

		"name": "address",

		"type": "string",

		"default": ""

		}

		]

		}

* Get the schema id from schema tab itself

  

* Update the schema id in the publisher and Run the application

* As we have created the new schema and it should have a schema id which we need to add in the publisher side.

* Start the application for testing: `go run main.go`

* The test main.go exposed a test rest-endpoint (http://localhost:8080/produce/message) and takes a employee object as an input text. Then it will **serialized the employee event using avro schema (get the avro schema from schema registry by its id and store in local cache) and convert the employee into serialized binary format** and then publish it into the `EMPLOYEE_AVRO_TOPIC`. Below is the sample employee text for rest endpoint:

  

			{

			"employeeName":"Monojit",

			"department":"Development",

			"salary":100.50,

			"address":"Bangalore"

			}

* To test if the publisher able to publish success fully or not, we have also started a consumer in the application as a separate go routine. This consumer will consume the event from the topic and then from its first 5 bytes try to parse the **schema id and based on the schema it will deserialized the event** and convert into text/json value.
