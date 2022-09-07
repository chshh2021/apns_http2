package apns_http2

type ReqData struct {
	Id     string
	Token  string
	Data   string
	Expire int64
}

type RespData struct {
	Id    string
	Token string
	Err   error
}

type Queue struct {
	reqChan  chan *ReqData
	Response chan *RespData
	client   *Client
}

func worker(q *Queue) {
	for it := range q.reqChan {
		err := q.client.Send(it.Token, it.Data, it.Expire)
		q.Response <- &RespData{it.Id, it.Token, err}
	}
}

func NewQueue(num int, client *Client) *Queue {
	q := &Queue{
		reqChan:  make(chan *ReqData, 1),
		Response: make(chan *RespData, 1),
		client:   client,
	}
	for i := 0; i < num; i++ {
		go worker(q)
	}

	return q
}

func (queue *Queue) Push(id, token, data string, expire int64) {
	queue.reqChan <- &ReqData{id, token, data, expire}
}

func (queue *Queue) Close() {
	close(queue.reqChan)

	close(queue.Response)
}
