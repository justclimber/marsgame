package helpers

import "time"

func makeTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeStampDif(t1, t2 time.Time) int64 {
	return makeTimestamp(t2) - makeTimestamp(t1)
}

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
