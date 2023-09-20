package event

type HandleResultConsumer interface {
	Consume(HandleResult)
}

type HandleBodyAndHeaders struct {
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
	Status  string              `json:"status"`
}

type HandleResult struct {
	ID string `json:"id"`

	URL       string                        `url:"method"`
	Method    string                        `json:"method"`
	Request   *HandleBodyAndHeaders         `json:"request"`
	Responses map[int]*HandleBodyAndHeaders `json:"responses"`
	Response  *HandleBodyAndHeaders         `json:"response"`
}
