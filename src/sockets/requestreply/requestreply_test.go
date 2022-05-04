package requestreply

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestRequestSocket(t *testing.T) {

	zmqContext, _ := zmq.NewContext()

	receiver, _ := zmqContext.NewSocket(zmq.REP)

	_ = receiver.Bind("tcp://*:9449")

	defer cleanupTest(receiver, zmqContext)

	var exchanges int = 10

	go StartRequestClient(exchanges)

	var requests []int64

	var replies []int64

	for i := 0; i < exchanges; i++ {

		message, _ := receiver.Recv(0)

		req, _ := strconv.Atoi(message)

		requests = append(requests, int64(req))

		time.Sleep(2 * time.Second)

		reply := rand.Int() + rand.Int()

		replies = append(replies, int64(reply))

		receiver.Send(strconv.Itoa(reply), 0)
	}

	receiver.Close()

	zmqContext.Term()

	assert.Equal(t, replies, serverReplies)
	assert.Equal(t, requests, clientRequests)
}

func TestReplySocket(t *testing.T) {

	var exchanges = 10

	go StartReplyServer(exchanges)

	zmqContext, _ := zmq.NewContext()

	sender, _ := zmqContext.NewSocket(zmq.REQ)

	_ = sender.Connect("tcp://localhost:9449")

	defer cleanupTest(sender, zmqContext)

	time.Sleep(1000 * time.Millisecond)

	var requests []int64

	var replies []int64

	for i := 0; i < exchanges; i++ {

		message := rand.Int() * rand.Intn(10)

		requests = append(requests, int64(message))

		sender.Send(strconv.Itoa(message), 0)

		time.Sleep(1 * time.Second)

		response, _ := sender.Recv(0)

		rep, _ := strconv.Atoi(response)

		replies = append(replies, int64(rep))
	}

	assert.Equal(t, replies, serverReplies)
	assert.Equal(t, requests, clientRequests)

}

func cleanupTest(socket *zmq.Socket, context *zmq.Context) {

	socket.Close()

	context.Term()
}
