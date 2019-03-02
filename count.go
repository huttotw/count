package count

import (
	"context"
	"sync"
	"time"
)

// Counter provides some abstraction to the Increment and Close
// functions. Normally, a new counter will be shared throughout
// a program. You can specify your own logger if you wish to see
// the internal logging of this package.
type Counter struct {
	Logger Logger

	m      map[string]uint64
	mutex  sync.Mutex
	window time.Duration
	w      Writer
	quit   chan struct{}
	errs   chan error
}

// New will create a new counter that will count keys and write
// the total for each window to the writer. The window will be the
// granularity you want to provide to the end users. The writer
// will receive the counts at every tick of the window.
func New(window time.Duration, w Writer) *Counter {
	c := Counter{
		Logger: noOpLogger{},

		m:      make(map[string]uint64),
		window: window,
		w:      w,
		quit:   make(chan struct{}),
		errs:   make(chan error),
	}
	go c.run()
	return &c
}

// Increment will increment the key by the delta. Increment is safe
// for concurrent use. After a write, increment will start at 0.
func (c *Counter) Increment(key string, delta uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.m[key] += delta
}

// Close should be called before exiting a program to ensure that the
// final, and perhaps incomplete window is sent to the writer.
func (c *Counter) Close() error {
	c.Logger.Log("msg", "closing counter", "keys", len(c.m))
	c.quit <- struct{}{}
	c.Logger.Log("msg", "closed counter")
	return nil
}

// write will send each key to the given writer in order to persist.
// After a key is written, it is deleted from the map. write is safe
// for concurrent use.
func (c *Counter) write(now time.Time) error {
	c.mutex.Lock()
	for k, v := range c.m {
		c.Logger.Log("msg", "writing count", "key", k, "value", v)
		err := c.w.Write(context.TODO(), k, v, now)
		if err != nil {
			return err
		}
		delete(c.m, k)
		c.Logger.Log("msg", "reset count", "key", k)
	}
	c.mutex.Unlock()
	return nil
}

// run will listen for each tick, error, or quit signal and handle it
// appropriately. run is a long running function and should be started
// when New is called. run will return after quit is called, and the
// counters have been written.
func (c *Counter) run() {
	tick := time.NewTicker(c.window).C
	for {
		select {
		case now := <-tick:
			c.Logger.Log("msg", "tick received", "time", now)
			err := c.write(now)
			if err != nil {
				c.errs <- err
				continue
			}
		case err := <-c.errs:
			c.Logger.Log("msg", "error received", "err", err.Error())
		case <-c.quit:
			c.Logger.Log("msg", "quit received")
			err := c.write(time.Now())
			if err != nil {
				c.errs <- err
				continue
			}
			c.errs <- nil
			return
		}
	}
}
