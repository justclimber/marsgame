package changelog

const ChannelSize = 10
const MaxLogSize = 8

type Log struct {
	ChangesByTimeCh  chan *ChangeByTime
	TerminateCh      chan bool
	ChangesByTimeLog []*ChangeByTime
}
