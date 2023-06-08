// https://www.cnblogs.com/my_life/articles/14870151.html

package util

import "time"

type QPS struct {
	intervals [2]Interval
	index     int
}

func New(interval, window int) QPS {
	bc := interval / window
	return QPS{
		intervals: [2]Interval{
			{
				buckets:     make([]Bucket, 0, bc),
				bucketCount: bc,
				window:      window,
			},
			{
				buckets:     make([]Bucket, 0, bc),
				bucketCount: bc,
				window:      window,
			},
		},
	}
}

func (q *QPS) Count() {
}

type Interval struct {
	buckets     []Bucket
	bucketCount int
	window      int
}

func (i *Interval) currentBucket() int {
	currentTime := time.Now().Nanosecond()
	return (currentTime / i.window) % i.bucketCount
}

type Bucket struct{}
