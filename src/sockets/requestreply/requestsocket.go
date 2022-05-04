package requestreply

import (
	zmq "github.com/pebbe/zmq4"
	"math/rand"
	"strconv"
	"time"
)

var clientRequests []int64

func StartRequestClient(exchanges int) {

	zmqContext, _ := zmq.NewContext()

	sender, _ := zmqContext.NewSocket(zmq.REQ)

	_ = sender.Connect("tcp://localhost:9449")

	defer cleanupClient(sender, zmqContext)

	time.Sleep(1000 * time.Millisecond)

	for i := 0; i < exchanges; i++ {

		message := rand.Int() * rand.Intn(10)

		clientRequests = append(clientRequests, int64(message))

		sender.Send(strconv.Itoa(message), 0)

		time.Sleep(1 * time.Second)

		response, _ := sender.Recv(0)

		rep, _ := strconv.Atoi(response)

		serverReplies = append(serverReplies, int64(rep))
	}

}

func cleanupClient(socket *zmq.Socket, context *zmq.Context) {

	socket.Close()
	context.Term()
}
