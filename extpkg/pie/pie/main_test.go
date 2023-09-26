package pie

import (
	"context"
	"os"
	"testing"
	"time"

	testify_stats "github.com/elliotchance/testify-stats"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(testify_stats.Run(m))
}

func assertImmutableStrings(t *testing.T, ss *Strings) func() {
	before := (*ss).JSONString()

	return func() {
		after := (*ss).JSONString()
		assert.Equal(t, before, after)
		assert.True(t, before == after)
	}
}

func assertImmutableInts(t *testing.T, ss *Ints) func() {
	before := (*ss).JSONString()

	return func() {
		after := (*ss).JSONString()
		assert.Equal(t, before, after)
		assert.True(t, before == after)
	}
}

func assertImmutableFloat64s(t *testing.T, ss *Float64s) func() {
	before := (*ss).JSONString()

	return func() {
		after := (*ss).JSONString()
		assert.Equal(t, before, after)
		assert.True(t, before == after)
	}
}

func createContextByDelay(t time.Duration) context.Context {
	ctx := context.Background()
	if t > 0 {
		ctx, _ = context.WithTimeout(ctx, t)
	}

	return ctx
}

func getFloat64sFromChan(ch chan float64, t time.Duration) func() Float64s {
	done := make(chan struct{})
	var c Float64s
	if t > 0 {
		go func() {
			ticker := time.NewTicker(t)
			defer ticker.Stop()
			for range ticker.C {
				val, ok := <-ch
				if !ok {
					break
				} else {
					c = append(c, val)
				}
			}
			done <- struct{}{}

		}()
	} else {
		go func() {
			for val := range ch {
				c = append(c, val)
			}
			done <- struct{}{}
		}()
	}

	return func() Float64s {
		<-done
		return c
	}
}

func getIntsFromChan(ch chan int, t time.Duration) func() Ints {
	done := make(chan struct{})
	var c Ints
	if t > 0 {
		go func() {
			ticker := time.NewTicker(t)
			defer ticker.Stop()
			for range ticker.C {
				val, ok := <-ch
				if !ok {
					break
				} else {
					c = append(c, val)
				}
			}
			done <- struct{}{}

		}()
	} else {
		go func() {
			for val := range ch {
				c = append(c, val)
			}
			done <- struct{}{}
		}()
	}

	return func() Ints {
		<-done
		return c
	}
}

func getStringsFromChan(ch chan string, t time.Duration) func() Strings {
	done := make(chan struct{})
	var c Strings
	if t > 0 {
		go func() {
			ticker := time.NewTicker(t)
			defer ticker.Stop()
			for range ticker.C {
				val, ok := <-ch
				if !ok {
					break
				} else {
					c = append(c, val)
				}
			}
			done <- struct{}{}

		}()
	} else {
		go func() {
			for val := range ch {
				c = append(c, val)
			}
			done <- struct{}{}
		}()
	}

	return func() Strings {
		<-done
		return c
	}
}
