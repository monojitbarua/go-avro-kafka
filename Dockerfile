FROM golang:1.17-alpine as build-env
 
# Set environment variable
ENV APP_NAME go-avro-kafka
ENV CMD_PATH main.go
 
# Copy application data into image
COPY . $GOPATH/src/$APP_NAME
COPY config.json $GOPATH/src/$APP_NAME/config.json
WORKDIR $GOPATH/src/$APP_NAME
 
# Budild application
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
 
# Run Stage
FROM alpine:3.14
 
# Set environment variable
ENV APP_NAME go-avro-kafka
 
# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
#COPY config.json /$APP_NAME/config.json 
 
# Expose application port
EXPOSE 8080
 
# Start app
CMD ./$APP_NAME