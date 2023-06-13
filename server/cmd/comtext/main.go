package main

import (
	"context"
	"fmt"
	"time"
)

type paramKey struct{}

func main() {
	c := context.WithValue(context.Background(), paramKey{}, "abc")
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	go mainTask(c)

	// 输入c主动cancel
	var cmd string
	for {
		fmt.Scan(&cmd)
		if cmd == "c" {
			cancel()
		}
	}
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with param %q\n", c.Value(paramKey{}))
	c1, cancel := context.WithTimeout(c, 2*time.Second) // 传c，c1相当于子任务
	defer cancel()
	// smallTask(context.Background(), "task3", 4*time.Second) // context.Background()为新的没有上面设置的5秒超时，即使完成时间设置成4秒，也会立刻打印出task3 done
	// smallTask(c, "task1", 1*time.Second) // 会打印出task1 done
	smallTask(c1, "task1", 4*time.Second) // c1超时时间设置成2秒，完成时间设置成4秒，会打印出task1 cancelled
	smallTask(c, "task2", 2*time.Second)  // c是同一个任务的步骤，c的超时时间是5秒，完成时间2秒，c1超时时间2秒还剩3秒, 会打印出task2 done

	// 启动后台任务
	go func() {
		c2, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 传context.Background()为后台任务
		defer cancel()
		smallTask(c2, "task4", 8*time.Second) // 新开的context.Background()相当于后台任务，新的后台任务不会携带参数，打印出task4 started with param %!q(<nil>)和 task4 done
	}()
}

func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
	select {
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		fmt.Printf("%s cancelled\n", name)
	}
}
