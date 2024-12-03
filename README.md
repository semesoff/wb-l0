# wb-l0

This is test task L0.

### Installation
```bash
git clone https://github.com/semesoff/wb-l0.git
cd wb_l0
```

### Launch
Launch containers: wb_l0_postgres, wb_l0_kafka.
```bash
docker-compose up
```
Create topic **orders** in kafka.
```bash
docker exec -it wb_l0_kafka /opt/kafka/bin/kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
```
Launch app. Producer automatically sends messages after launch.
```bash
go run cmd/app/main.go
```

### Usage
```bash
http://localhost:8000/orders/b563feb7b2b84b6test
```
