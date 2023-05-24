package timingwheel

import (
	"errors"
	"github.com/oldbai555/lbtool/pkg/timingwheel/delayqueue"
	"sync/atomic"
	"time"
	"unsafe"
)

// TimingWheel is an implementation of Hierarchical Timing Wheels.
// 分层时间轮
type TimingWheel struct {
	tick      int64 // 每一个时间格的跨度,以毫秒为单位
	wheelSize int64 // 时间格的数量

	interval    int64                  // 总的跨度数 tick * wheelSize，以毫秒为单位
	currentTime int64                  // 当前指针指向的时间，以毫秒为单位
	buckets     []*bucket              // 时间格列表
	queue       *delayqueue.DelayQueue // 延迟队列,

	// The higher-level overflow wheel.
	//
	// NOTE: This field may be updated and read concurrently, through Add().
	overflowWheel unsafe.Pointer // 上一层时间轮的指针

	exitC     chan struct{} // 退出通知
	waitGroup waitGroupWrapper
}

// NewTimingWheel creates an instance of TimingWheel with the given tick and wheelSize.
// 对外暴露的初始化时间轮方法,参数为时间格跨度，和时间格数量
func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	// 时间格(毫秒)
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be greater than or equal to 1ms"))
	}
	// 开始时间
	startMs := timeToMs(time.Now().UTC())

	return newTimingWheel(
		tickMs,
		wheelSize,
		startMs,
		delayqueue.New(int(wheelSize)), // 延时队列
	)
}

// newTimingWheel is an internal helper function that really creates an instance of TimingWheel.
// 内部初始化时间轮的方法，参数为，时间格跨度（毫秒），时间格数量，开始时间（毫秒）, 延迟队列
func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, queue *delayqueue.DelayQueue) *TimingWheel {
	// 根据时间格数量创建时间格列表
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	return &TimingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * wheelSize,
		buckets:     buckets,
		queue:       queue,
		exitC:       make(chan struct{}),
	}
}

// add inserts the timer t into the current timing wheel.
// add添加定时器到时间轮
// 如果定时器已过期返回false
func (tw *TimingWheel) add(t *Timer) bool {
	// 当前秒钟时间轮的currentTime = 1626333377000（2021-07-15 15:16:17）
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if t.expiration < currentTime+tw.tick {
		// Already expired
		// 定时器的过期时间已经过期返回false
		return false
	} else if t.expiration < currentTime+tw.interval {
		// 定时器的过期时间小于当前时间轮的当前时间+轮的总跨度，将定时器放到对应的bucket中，并将bucket放入延迟队列。

		// 假设过期时间为2021-07-15 15:17:02（1626333422000）
		// 1626333422000 < 1626333377000 + 60*1000
		// virtualID = 1626333422000 / 1000 = 1626333422
		// 1626333422%60 = 2，将定时器放到第2个时间格中
		// 设置bucket（时间格）的过期时间

		// Put it into its own bucket
		virtualID := t.expiration / tw.tick
		b := tw.buckets[virtualID%tw.wheelSize]
		b.Add(t)

		// Set the bucket expiration time
		if b.SetExpiration(virtualID * tw.tick) {
			// The bucket needs to be enqueued since it was an expired bucket.
			// We only need to enqueue the bucket when its expiration time has changed,
			// i.e. the wheel has advanced and this bucket get reused with a new expiration.
			// Any further calls to set the expiration within the same wheel cycle will
			// pass in the same value and hence return false, thus the bucket with the
			// same expiration will not be enqueued multiple times.

			// 如果设置的过期时间不等于桶的过期时间
			// 将bucket添加到延迟队列，重新排序延迟队列

			tw.queue.Offer(b, b.Expiration())
		}

		return true
	} else {
		// Out of the interval. Put it into the overflow wheel
		// 定时器的过期时间 大于 当前时间轮的当前时间+轮的总跨度，递归将定时器添加到上一层轮。
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel == nil {
			atomic.CompareAndSwapPointer(
				&tw.overflowWheel,
				nil,
				unsafe.Pointer(newTimingWheel(
					tw.interval,
					tw.wheelSize,
					currentTime,
					tw.queue,
				)),
			)
			overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
		}
		return (*TimingWheel)(overflowWheel).add(t)
	}
}

// addOrRun inserts the timer t into the current timing wheel, or run the
// timer's task if it has already expired.
// 执行已过期定时器的任务，将未到期的定时器重新放回时间轮
func (tw *TimingWheel) addOrRun(t *Timer) {
	if !tw.add(t) {
		// Already expired

		// Like the standard time.AfterFunc (https://golang.org/pkg/time/#AfterFunc),
		// always execute the timer's task in its own goroutine.
		go t.task()
	}
}

// 推进时钟
// 我们就以时钟举例:假如当前时间是2021-07-15 15:16:17（1626333375000毫秒），过期时间是2021-07-15 15:17:18（1626333438000）毫秒
// 从秒轮开始,1626333438000 > 1626333377000 + 1000， truncate(1626333438000,1000)=1626333438000, 秒轮的当前时间设置为1626333438000（2021-07-15 15:17:18），有上层时间轮
// 到了分钟轮 1626333438000 > 1626333377000 + 60000=1626333437000, truncate(1626333438000,60000)=1626333420000（2021-07-15 15:17:00），分轮的当前时间设置为1626333438000，有上层时间轮
// 到了时钟轮 1626333438000 < 1626333377000 + 360000，时钟轮当前时间不变（2021-07-15 15:16:17），没上层时间轮
func (tw *TimingWheel) advanceClock(expiration int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= currentTime+tw.tick {

		// 将过期时间截取到时间格间隔的最小整数倍
		// 举例：
		// expiration = 100ms，tw.tick = 3ms, 结果 100 - 100%3 = 99ms,因此当前的时间来到了99ms，
		// 目的就是找到合适的范围，比如[0,3)、[3-6)、[6,9) expiration=5ms时，currentTime=3ms。

		currentTime = truncate(expiration, tw.tick)
		atomic.StoreInt64(&tw.currentTime, currentTime)

		// Try to advance the clock of the overflow wheel if present
		// 如果有上层时间轮，那么递归调用上层时间轮的引用
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advanceClock(currentTime)
		}
	}
}

// Start starts the current timing wheel.
// 时间轮转起来
func (tw *TimingWheel) Start() {
	tw.waitGroup.Wrap(func() {
		// 开启一个协程，死循环延迟队列，将已过期的bucket(时间格)弹出
		tw.queue.Poll(tw.exitC, func() int64 {
			return timeToMs(time.Now().UTC())
		})
	})

	tw.waitGroup.Wrap(func() {
		for {
			select {
			// 开启另外一个协程，阻塞接收延迟队列弹出的bucket(时间格)
			case elem := <-tw.queue.C:
				// 从延迟队列弹出来的是一个bucket（时间格）
				b := elem.(*bucket)
				// 时钟推进，将时钟的当前时间推进到过期时间
				tw.advanceClock(b.Expiration())
				// 将bucket（时间格）中的已到期的定时器执行，还没有到过期时间重新放回时间轮
				b.Flush(tw.addOrRun)
			case <-tw.exitC:
				return
			}
		}
	})
}

// Stop stops the current timing wheel.
//
// If there is any timer's task being running in its own goroutine, Stop does
// not wait for the task to complete before returning. If the caller needs to
// know whether the task is completed, it must coordinate with the task explicitly.
// 停止时间轮
// 关闭管道
func (tw *TimingWheel) Stop() {
	close(tw.exitC)
	tw.waitGroup.Wait()
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
// 添加定时任务到时间轮
func (tw *TimingWheel) AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{
		expiration: timeToMs(time.Now().UTC().Add(d)),
		task:       f,
	}
	tw.addOrRun(t)
	return t
}

// Scheduler determines the execution plan of a task.
type Scheduler interface {
	// Next returns the next execution time after the given (previous) time.
	// It will return a zero time if no next time is scheduled.
	//
	// All times must be UTC.
	Next(time.Time) time.Time
}

// ScheduleFunc calls f (in its own goroutine) according to the execution
// plan scheduled by s. It returns a Timer that can be used to cancel the
// call using its Stop method.
//
// If the caller want to terminate the execution plan halfway, it must
// stop the timer and ensure that the timer is stopped actually, since in
// the current implementation, there is a gap between the expiring and the
// restarting of the timer. The wait time for ensuring is short since the
// gap is very small.
//
// Internally, ScheduleFunc will ask the first execution time (by calling
// s.Next()) initially, and create a timer if the execution time is non-zero.
// Afterwards, it will ask the next execution time each time f is about to
// be executed, and f will be called at the next execution time if the time
// is non-zero.
func (tw *TimingWheel) ScheduleFunc(s Scheduler, f func()) (t *Timer) {
	expiration := s.Next(time.Now().UTC())
	if expiration.IsZero() {
		// No time is scheduled, return nil.
		return
	}

	t = &Timer{
		expiration: timeToMs(expiration),
		task: func() {
			// Schedule the task to execute at the next time if possible.
			expiration := s.Next(msToTime(t.expiration))
			if !expiration.IsZero() {
				t.expiration = timeToMs(expiration)
				tw.addOrRun(t)
			}

			// Actually execute the task.
			f()
		},
	}
	tw.addOrRun(t)

	return
}
