package storage

type item struct {
	heartbeatCount int
	checkCount     int
}

func (i *item) IncrementHeartbeatCount() {
	i.heartbeatCount++
}

func (i *item) ResetHeartbeatCounter() {
	i.heartbeatCount = 0
}

func (i *item) IncrementCheckCount() {
	i.checkCount++
}

func (i *item) ResetCheckCounter() {
	i.checkCount = 0
}

func newItem() *item {
	return &item{
		heartbeatCount: 0,
		checkCount:     0,
	}
}
