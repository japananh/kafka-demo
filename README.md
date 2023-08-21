# How to run a simple Kafka Cluster?

## Run Kafka with Confluent and Go

Full tutorials (in Vietnamese) can be found on: https://200lab.io/blog/

To run Kafka producer

```bash
cd kafka-using-confluent

# Run Kafka producer
go run producer.go \
<bootstrap-servers> \
<cluster-api-key> \
<cluster-api-secret> \
<topic> \
<schema-registry-url> \
<schema-registry-api-key> \
<schema-registry-api-secret>
```

To run Kafka consumer

```bash
# Run Kafka consumer
go run consumer.go \
<bootstrap-servers> \
<cluster-api-key> \
<cluster-api-secret> \
<consumer-group-id> \
<topic-1> \
<topic-2> \
...
<topic-N>
```

## Run Kafka using Docker Compose

```bash
cd kafka-using-docker-compose
docker comopose up -d
```