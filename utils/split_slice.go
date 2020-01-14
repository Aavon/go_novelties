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

func (iter *SplitIter) Next() (start, end int, ok, hasNext bool) {
	finished := atomic.LoadInt32(&iter.finish) == 1
	if finished {
		return int(iter.nextStart), int(iter.nextEnd), true, false
	}
	oldEnd := atomic.LoadInt64(&iter.nextEnd)
	if !atomic.CompareAndSwapInt64(&iter.nextStart, iter.nextStart, oldEnd) {
		return
	}
	newEnd := oldEnd + iter.step
	if newEnd >= iter.l {
		newEnd = iter.l
		if !atomic.CompareAndSwapInt32(&iter.finish, iter.finish, 1) {
			return
		}
		finished = true
	}
	if !atomic.CompareAndSwapInt64(&iter.nextEnd, iter.nextEnd, newEnd) {
		return
	}
	return int(oldEnd), int(newEnd), true, !finished
}
