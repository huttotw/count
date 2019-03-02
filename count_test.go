package count

import (
	"context"
	"testing"
	"time"
)

type mockWriter struct {
	writes int
}

func (w *mockWriter) Write(ctx context.Context, key string, value uint64, current time.Time) error {
	w.writes++
	return nil
}

func TestCounter_run(t *testing.T) {
	w := mockWriter{}
	c := New(time.Second, &w)
	c.m = map[string]uint64{
		"unit-test": 1,
	}

	time.Sleep(time.Second * 2)

	if w.writes != 1 {
		t.Fatal("expected one write")
	}
}

func TestCounter_write(t *testing.T) {
	w := mockWriter{}
	c := New(time.Second, &w)
	c.m = map[string]uint64{
		"unit-test": 1,
	}

	now := time.Now()
	err := c.write(now)
	if err != nil {
		t.Fatal(err)
	}

	if w.writes != 1 {
		t.Fatal("expected one write")
	}
	if len(c.m) != 0 {
		t.Fatal("expected map to be empty after write")
	}
}

func TestCounter_Increment(t *testing.T) {
	w := mockWriter{}
	c := New(time.Second, &w)

	c.Increment("unit-test", 1)
	if c.m["unit-test"] != 1 {
		t.Fatal("expected value to be counted")
	}

	c.Increment("unit-test", 2)
	if c.m["unit-test"] != 3 {
		t.Fatal("expected value to be counted")
	}
}

func BenchmarkCounter_Increment(b *testing.B) {
	w := mockWriter{}
	c := New(time.Second, &w)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Increment("unit-test", 1)
	}
}

func TestCounter_Close(t *testing.T) {
	w := mockWriter{}
	c := New(time.Second, &w)
	c.m = map[string]uint64{
		"unit-test": 1,
	}

	err := c.Close()
	if err != nil {
		t.Fatal(err)
	}

	if w.writes != 1 {
		t.Fatal("expected to flush map")
	}
	if len(c.m) != 0 {
		t.Fatal("expected map to be empty")
	}
}