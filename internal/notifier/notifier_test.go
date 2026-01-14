package notifier

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/assert"
)

func TestNotifier(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		n := &Notifier{}

		var x []int
		go func() {
			ch, cancel := n.Subscribe()
			defer cancel()

			<-ch
			x = append(x, 1)
		}()

		go func() {
			ch, cancel := n.Subscribe()
			defer cancel()

			<-ch
			x = append(x, 2)
		}()

		// Make sure the goroutines are durable blocked on their channel receives first
		synctest.Wait()

		n.Notify()
		synctest.Wait()

		slices.Sort(x)
		assert.Equal(t, []int{1, 2}, x)

		// Should be able to notify after the subscriptions are done listening, nothing happens
		n.Notify()

		// Test receiving multiple notifications
		x = x[:0]
		go func() {
			ch, cancel := n.Subscribe()
			defer cancel()

			for {
				select {
				case <-t.Context().Done():
					return
				case <-ch:
					x = append(x, 1)
				}
			}
		}()

		synctest.Wait()

		n.Notify()
		synctest.Wait()

		slices.Sort(x)
		assert.Equal(t, []int{1}, x)

		go func() {
			ch, cancel := n.Subscribe()
			defer cancel()

			for {
				select {
				case <-t.Context().Done():
					return
				case <-ch:
					x = append(x, 2)
				}
			}
		}()

		synctest.Wait()

		n.Notify()
		synctest.Wait()

		slices.Sort(x)
		assert.Equal(t, []int{1, 1, 2}, x)

		n.Notify()
		synctest.Wait()

		slices.Sort(x)
		assert.Equal(t, []int{1, 1, 1, 2, 2}, x)
	})
}
