# stream-cardinality-excercise

Repo that represents my first interaction with kafka and streams and adventures to the world of Probabilistic algorithms - in this case measuring cardinality of an open set (streaming data).

## Requirements

- Install [golang](https://golang.org/doc/install)
- Setup [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)
- Install [dep](https://github.com/golang/dep)

## Buiding

- Clone this repository
- `dep ensure`
- `cd cmd`
- `go build -o excercise`

There is a Makefile available with `make all`. If you select this approach the binaries will be created in `./.gopath~/bin/`. Not that you need to run `dep ensure` prior `make all`.

## Running

- Run zookeper and kafka:

    ```bash
    brew services start zookeeper
    brew services start kafka
    ```
- Create a topic:

    ```bash
    kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic users
    ```
- Start console producer:

    ```bash
    kafka-console-producer --broker-list localhost:9092 --topic 'users' < ~/code/data/stream.jsonl
    ```

Before running just note that here processor is equivalent to consumer. I am using the slang of the library I am using
to provide information.

- To run processor that just outputs the messages from kafka topic run

    ```bash
    # Using go run
    go run cmd/* processor-stdout --address localhost:9092 --topic users
    # Using binary
    ./.gopath~/bin/stream-cardinality-excercise processor-stdout --address localhost:9092 --topic users
    ```

- To run processor that can count by frame-per-rate method:

    ```bash
    # Using go run
    go run cmd/* processor-counter --address localhost:9092 --topic users --rate 1s --output-rate 5s --json 1
    # Using binary
    ./.gopath~/bin/stream-cardinality-excercise processor-counter --address localhost:9092 --topic users --rate 1s --output-rate 5s --json 1

    ```

- To run processor that counts unique users per minute run:

```bash
# Using go run
go run cmd/* processor-user --address localhost:9092 --topic users

# Using binary
./.gopath~/bin/stream-cardinality-excercise processor-user --address localhost:9092
```

## Notes, Benchmarks, Thoughts

Available [here](doc/notes.md)