package upstream

import "strconv"

type UpstreamID string

type Node struct {
	Host    string
	Port    int
	Weight  int
	Healthy bool
	Tags    map[string]string
}

type Upstream struct {
	ID    UpstreamID
	Name  string
	Nodes []Node
}

func (n Node) URL() string {
	return "http://" + n.Host + ":" + strconv.Itoa(n.Port)
}
