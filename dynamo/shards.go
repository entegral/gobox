package dynamo

type Shard struct {
	Base     string
	maxShard int
}

func (s Shard) MaxShard() int {
	if s.maxShard == 0 {
		return 100
	}
	return s.maxShard
}

func (s Shard) SetMaxShard(max int) {
	s.maxShard = max
}
