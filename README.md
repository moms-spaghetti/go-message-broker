# running the broker and client
The broker is in `cmd/broker/main.go`  
The client is in `cmd/client/main.go`  
Use `go run main.go` in each folder to run both

# broker
## support
- TCP / UDP
- Creating topics (map of topic id string and subscribers)
- Creating subscribers
- Getting topics & subscribers
- Getting only subscribers
- Subscribing to a topic
- Downloading subscriber messages

## subscriber messages
Subscriber messages are held on subscriber objects as a map of string (topic id) messages. Messages are objects containing id, topic, body and a done flag

## downloading messages
Downloading subscriber messages is continuous.  
Messages will be downloaded FIFO until there are no further messages.

## working on messages
Once downloaded, work is simulated using `time.Sleep`. After, a response is sent to the broker flagging the message as done (bool). It's removed from the subscriber's message list for the topic. If work is cut short the message will not be removed from the subscriber message list and will be re-sent upon the next download attempt.

## notes
Topics are basic maps of string ids and subscribers  
Subscribers are a map of string ids (subscriber id) and subscriber objects  
All are held on the broker

# client app
A basic client app exists to improve communication with the broker. Runs in terminal  
Supports TCP and UDP functions

# client app options
## menu 1
- TCP
- UDP
- QUIT

## menu 2
### publisher options
- add topic
- get topics and subscribers
- send message to topic/subscribers (takes the topicid and a body of any type)
	
### broker options
- create new subscriber (returns a subscriber object with id)
- list subscribers

### subscriber options
- subscribe to topic (takes a subscriber id and topic id)
- download subscriber messages (starts the process of downloading subscriber messages)

### other options
- QUIT

