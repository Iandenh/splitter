package upstream

type Upstream struct {
	Url string `url:"url"`
}

var (
	upstreams []Upstream
)

func AddUpstream(u Upstream) {
	upstreams = append(upstreams, u)
}

func GetUpstreams() []Upstream {
	return upstreams
}
