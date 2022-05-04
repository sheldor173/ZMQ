package requestreply

import (
	zmq "github.com/pebbe/zmq4"
	"math/rand"
	"strconv"
	"time"
)

var serverReplies []int64

func StartReplyServer(exchanges int) {

	zmqContext, _ := zmq.NewContext()

	receiver, _ := zmqContext.NewSocket(zmq.REP)

	_ = receiver.Bind("tcp://*:9449")

	defer cleanupServer(receiver, zmqContext)

	for i := 0; i < exchanges; i++ {

		message, _ := receiver.Recv(0)

		req, _ := strconv.Atoi(message)

		clientRequests = append(clientRequests, int64(req))

		time.Sleep(2 * time.Second)

		reply := rand.Int() + rand.Int()

		serverReplies = append(serverReplies, int64(reply))

		receiver.Send(strconv.Itoa(reply), 0)
	}
}

func cleanupServer(socket *zmq.Socket, context *zmq.Context) {

	socket.Close()
	context.Term()
}
