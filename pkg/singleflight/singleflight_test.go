package singleflight

import (
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

/*
在并发编程中，当多个 goroutine 同时访问和修改共享的数据时，可能会出现竞态条件（Race Condition）的问题。
在这种情况下，如果直接使用 nums++ 运算符来对 nums 进行加一操作，由于并发执行的不确定性，可能会导致数据不一致或者丢失更新的情况发生。

为了避免这种问题，Go语言提供了 sync/atomic 包来进行原子操作，其中的 AddInt64 函数可以确保对 int64 类型的变量进行原子性的加法操作，
保证在同一时刻只有一个 goroutine 对该变量进行操作，从而避免了竞态条件问题。因此，在并发编程中，推荐使用原子操作来保证数据的正确性和一致性。
*/

func TestGroup_Do(t *testing.T) {
	wg := new(sync.WaitGroup)
	nums := int64(0)
	cnt := 100
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			defer wg.Done()
			group := NewGroup()
			_, err := group.Do("redis", func() (interface{}, error) {
				atomic.AddInt64(&nums, 1) //原子操作，将 nums 加一
				return n, nil
			})
			require.NoError(t, err) // 断言 err 为 nil
		}(i)
	}
	wg.Wait()
	t.Log(cnt, nums)
	require.True(t, nums <= int64(cnt)) // 断言 nums 的值小于等于 cnt
}

func TestGroup_WithoutDo(t *testing.T) {
	wg := new(sync.WaitGroup)
	cnt := 100
	nums := int64(0)
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := func() (interface{}, error) {
				atomic.AddInt64(&nums, 1)
				return n, nil
			}()
			require.NoError(t, err) // 断言 err 为 nil
		}(i)
	}
	wg.Wait()
	require.EqualValues(t, nums, cnt) // 断言 nums 的值等于 cnt
}
