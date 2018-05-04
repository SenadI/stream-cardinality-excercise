# Notes/Benchmars/Thoughts

## 1st part

- First check [kafka](https://kafka.apache.org/documentation/#gettingStarted).

- Read about Kafka architecture and try to identify the problem within the kafka documentation. In general it is clear even from the problem readme that this is stream processing - Incoming topic messages are processed and the outputs are sent to another topic.

- I will need a consumer to consume the messages, counting mechanism, and then output messages to other kafka-topic.

- Lets install kafka first. I use OsX at home.

    ```bash
    brew install kafka
    ==> Pouring kafka-1.1.0.high_sierra.bottle.tar.gz
    ==> Caveats
    To have launchd start kafka now and restart at login:
    brew services start kafka
    Or, if you don't want/need a background service you can just run:
    zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
    ==> Summary
    🍺  /usr/local/Cellar/kafka/1.1.0: 157 files, 47.8MB
    ```

- Ok so I've heard about zookeeper but didn't use it. Google, read [docs](https://zookeeper.apache.org/) - Get it: centralized service for dist sync, group services, naming and cfg.

- Language - I'll do GO the requirement is too use the fastest in your arsenal

- First research the counting. I know this is cardinality issue - problem is that data is coming in streams.

- A bit of googleing (Big Data Cardinality/Counters etc.) on cardinality counters and streams. Its seams there are several algorithms, the most google results come for HyperLogLog and HyperLogLog+. I visited another [site](http://highscalability.com/blog/2012/4/5/big-data-counting-how-to-count-a-billion-distinct-objects-us.html) where there are compairisons done on 3 algorithms.

    | Counter        | Bytes           | Count  |
    | ------------- |:-------------:| -----:|
    | Hash Set      | 10447016 | 67801 |
    | Linear      | 3384      |   67080 |
    | HyperLogLog | 512      |    70002 |

    Ok so HyperLogLog seams to be our man but its too complicated to implement for this time - so I reserach golang HyperLogLog implementation and I gind suitable implementation and I I decide to implement the LinarProbabilistic counter and port the HyperLogLog from the reference implementation - Linear Probablistic counter is easy to implement and would serve as a good learning problem.

- Back to kafka - I Lets write a simple consumer that consumes the messages and outputs the metric (count total messages) and  (#messages consumed / second ). This will get me started. I check out the reference implementations and find a suitable [match](https://github.com/movio/kasper)

- After some time = Useful commands

    ```bash
    # restart zookeper
    brew services restart zookeeper
    # restart kafka
    brew services restart kafka
    # form kafka webpage on how to start a topic
    kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic "users"
    # console producer can read os.Stdin
    kafka-console-producer --broker-list localhost:9092 --topic 'users' < ~/code/data/stream.jsonl
    ```

## 2nd Part

- Created a simple stdout processor that outputs the messages on stdout.
    ```bash
    # Terminal 1
    github.com/senadi/stream-cardinality-excercise via 🐹 v1.10.2
    ➜ kafka-console-producer --broker-list localhost:9092 --topic 'users'
    >test
    # Terminal 2
    ➜ go run cmd/* processor-stdout --address localhost:9092 --topic 'users'
    {"level":"info","msg":"Topic processor is running..."}
    {"key":"","level":"info","msg":"Message received","offset":115679,"partition":0,"topic":"users","value":"test"}
    ```
- Lets load the data and count how many frames per rate is our consumer consuming. There is a great library I used before for rate counting in general I will use it [here](https://github.com/paulbellamy/ratecounter), it can allow us to count by any duration frames per sec/min/h/d etc.

    ```bash
    # Terminal 1
    ➜ kafka-console-producer --broker-list localhost:9092 --topic 'users' < ~/code/data/stream.jsonl
    # Terminal 2
    ➜ go run cmd/* processor-counter --address localhost:9092 --topic users --rate 1s --output-rate 2s --json 0
    "frames":0,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:53:52+02:00"}
    {"level":"info","msg":"Topic processor is running...","time":"2018-05-04T14:53:52+02:00"}
    {"frames":0,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:53:54+02:00"}
    {"frames":9000,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:53:56+02:00"}
    {"frames":23000,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:53:58+02:00"}
    {"frames":25000,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:54:00+02:00"}
    {"frames":25616,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:54:02+02:00"}
    ```

The problem here with the data that is provided is that its JSON, json decoding is generally expensive - there are other posible formats I could suggest ( maybe even NCSA ) but JSON is generally a must due to largly multi-level nested data. I will edit the processor-counter and a struct option wheather do decode JSON or not.

Also extracted the rate (request-per-rate) and the results rate (when to output results lets say every 2 seconds).So lets measure request-per-rate for counter that desiralizes json.

With json decoding:

```bash
# Terminal 1
➜ kafka-console-producer --broker-list localhost:9092 --topic 'users' < ~/code/data/stream.jsonl
# Terminal 2
➜ go run cmd/* processor-counter --address localhost:9092 --topic users --rate 1s --output-rate 2s --json 0
{"frames":0,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:02+02:00"}
{"frames":910,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:04+02:00"}
{"frames":11851,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:10+02:00"}
{"frames":11033,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:12+02:00"}
{"frames":11685,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:14+02:00"}
{"frames":12704,"level":"info","msg":"Number of frames per duration","time":"2018-05-04T14:56:18+02:00"}
```

Its obvious to see that json decoding introduces overhead and affects the counting drastically.

## 3rd part

- For the timeline counter I will need first to fastly extract the `uid` and `ts` details from the messages. A bit of research and I found [gjson](https://github.com/tidwall/gjson).

- I will assume that messages are ordered - this simplifies the problem a bit.

The implementation turned out to be easy but at first it was challenging to wrap my head around it. One thing here again
this turned out to be easy becuase its for fixed 60seconds period with assumption that messages are ordered. This would 
most definitively turn into a much smarter structure/algorithm if the assumption is taken away.

```go
import (
	"hash/fnv"
)

func getHash(value string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(value))
	return h.Sum32()
}

// Counter that stors unique users per minute. The main structure is a HyperLogLog
// counter and the starting epoh time - The first incoming message. This implementation
// assumes messages are ordered.
type Counter struct {
	Hll *HyperLogLog
	// Now - artifical now of the incoming data reprents the first timestamp of the open interval.
	Now int
}

// Count unique uids from one minute timeframe.
func (c *Counter) Count(uid string, ts int) (bool, int) {
	// this is helpful for first run. Counter is started from 0 time
	// which is a signal that this is a first iteration.
	if c.Now == 0 {
		c.Now = ts
		c.Hll = NewHyperLogLog(0.1)
		c.Hll.Add(getHash(uid))
		return false, 0
		// Now if the ts is in the current minute range
		// Increment the counter
	} else if c.Now <= ts && c.Now+60 > ts {
		c.Hll.Add(getHash(uid))
		return false, 0
		// Now new one minute timeframe is here so report the curent count
	} else {
		uniqueCount := c.Hll.Count()
		c.Now = ts
		c.Hll = NewHyperLogLog(0.01)
		return true, int(uniqueCount)
	}
}
```

Example of the algorithm output:

```bash
# Terminal 1
➜ kafka-console-producer --broker-list localhost:9092 --topic 'users' < ~/code/data/stream.jsonl
# Terminal 2
➜ go run cmd/* processor-user --address localhost:9092 --topic users
{"level":"info","msg":"Topic processor is running...","time":"2018-05-04T17:47:35+02:00"}
{"count":250,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:40+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:43+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:46+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:48+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:50+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:53+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:55+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T17:47:58+02:00"}
{"count":32340,"level":"info","msg":"Counting...","time":"2018-05-04T17:48:01+02:00"}
{"count":30015,"level":"info","msg":"Counting...","time":"2018-05-04T17:48:03+02:00"}
{"count":31328,"level":"info","msg":"Counting...","time":"2018-05-04T17:48:08+02:00"}
{"count":31822,"level":"info","msg":"Counting...","time":"2018-05-04T17:48:13+02:00"}
```

## Project benchmarks

Among others already mentioned above (rate counting, json processing) I Decided to use `time` to measure the user processor. I had to modify the user processor from infinite loop to 5*10^6 messages. So I am going to time how much time do I need to
count unique users per minute: 

```bash
time go run cmd/* processor-user --address localhost:9092 --topic users
➜ time go run cmd/* processor-user --address localhost:9092 --topic users
(KASPER) 2018/05/04 18:53:37 INFO Consuming topic partition users-0 from offset '14905045' (newest offset is '16402086')
(KASPER) 2018/05/04 18:53:37 INFO Entering run loop
{"level":"info","msg":"Topic processor is running...","time":"2018-05-04T18:53:37+02:00"}
{"count":250,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:38+02:00"}
{"count":30090,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:40+02:00"}
{"count":28466,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:41+02:00"}
{"count":32704,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:42+02:00"}
{"count":31781,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:43+02:00"}
{"count":34474,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:45+02:00"}
{"count":30626,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:46+02:00"}
{"count":31822,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:48+02:00"}
{"count":32318,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:49+02:00"}
{"count":31240,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:50+02:00"}
{"count":31409,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:51+02:00"}
{"count":32645,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:53+02:00"}
{"count":34417,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:54+02:00"}
{"count":29742,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:55+02:00"}
{"count":31200,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:57+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T18:53:58+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:00+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:01+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:02+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:05+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:07+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:10+02:00"}
{"count":32340,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:12+02:00"}
{"count":30015,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:14+02:00"}
{"count":14487,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:15+02:00"}
{"count":31200,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:17+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:20+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:23+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:25+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:27+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:31+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:33+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:35+02:00"}
{"count":32340,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:36+02:00"}
{"count":30015,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:37+02:00"}
{"count":31328,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:39+02:00"}
{"count":31822,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:40+02:00"}
{"count":33650,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:42+02:00"}
{"count":31396,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:43+02:00"}
{"count":34165,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:45+02:00"}
{"count":31335,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:46+02:00"}
{"count":14441,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:47+02:00"}
{"count":31200,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:48+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:49+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:51+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:52+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:54+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T18:54:57+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:00+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:03+02:00"}
{"count":32340,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:05+02:00"}
{"count":30015,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:06+02:00"}
{"count":31328,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:07+02:00"}
{"count":31822,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:09+02:00"}
{"count":33650,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:15+02:00"}
{"count":31396,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:17+02:00"}
{"count":34165,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:20+02:00"}
{"count":31335,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:23+02:00"}
{"count":14441,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:30+02:00"}
{"count":31200,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:33+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:35+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:39+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:41+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:44+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:47+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:50+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:53+02:00"}
{"count":32340,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:56+02:00"}
{"count":30015,"level":"info","msg":"Counting...","time":"2018-05-04T18:55:59+02:00"}
{"count":31328,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:01+02:00"}
{"count":31822,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:04+02:00"}
{"count":33650,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:07+02:00"}
{"count":31396,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:10+02:00"}
{"count":34165,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:13+02:00"}
{"count":31335,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:16+02:00"}
{"count":14441,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:21+02:00"}
{"count":31200,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:24+02:00"}
{"count":29724,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:27+02:00"}
{"count":34999,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:31+02:00"}
{"count":32528,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:35+02:00"}
{"count":30254,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:38+02:00"}
{"count":32061,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:41+02:00"}
{"count":32682,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:45+02:00"}
{"count":33503,"level":"info","msg":"Counting...","time":"2018-05-04T18:56:49+02:00"}
go run cmd/* processor-user --address localhost:9092 --topic users

# RESULTS
80.82s user 34.82s system 59% cpu 3:15.82 total
....

```

## TODO

- Modify existing counter to accept per minute, per h. Its easy I just have one fixed 60s the change its minimal but no time left for the excercise.
- Create counter that can handle non-ordered messages.
- Read more about kafka - I wish I took one day just to play with kafka and then do this. :)
- This does not scale it is all in memory -  We can shurly use partition data and have consumers per partition. Time performance will be better because we can do in parallel.
- Write tests for the project.