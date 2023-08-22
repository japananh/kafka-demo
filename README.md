# How to run a simple Kafka Cluster?

## Run Kafka with Confluent and Go

Full Vietnamese instruction can be found on: [Hướng dẫn sử dụng Kafka với Confluent và Go](https://200lab.io/blog/huong-dan-su-dung-kafka/).

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