package utils

import "sync/atomic"

type SplitIter struct {
	l         int64
	nextStart int64
	nextEnd   int64
	step      int64
	finish    int32
}

func NewSplitIter(length int, step int) *SplitIter {
	return &SplitIter{
		l:         int64(length),
		step:      int64(step),
		nextStart: 0,
		nextEnd:   0,
	}
}

func (iter *SplitIter) Next() (start, end int, hasNext bool) {
	finished := atomic.LoadInt32(&iter.finish) == 1
	if finished {
		return int(iter.nextStart), int(iter.nextEnd), false
	}
	oldEnd := atomic.LoadInt64(&iter.nextEnd)
	atomic.StoreInt64(&iter.nextStart, oldEnd)
	newEnd := oldEnd + iter.step
	if newEnd >= iter.l {
		newEnd = iter.l
		if !atomic.CompareAndSwapInt32(&iter.finish, iter.finish, 1) {
			return 0, 0, false
		}
		finished = true
	}
	atomic.StoreInt64(&iter.nextEnd, newEnd)
	return int(oldEnd), int(newEnd), !finished
}
