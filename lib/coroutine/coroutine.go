package coroutine

import (
	log "lib/logwrap"
)

//注册协程函数，并执行
func Start(f func(current int, total int), total int) error {
	if total < 1 {
		return nil
	}

	chs := make([]chan int, total)
	for i := 0; i < total; i++ {
		chs[i] = make(chan int)
		current := i
		go coroutine_proxy(chs[i], current, total, f)
	}
	for _, ch := range chs {
		<-ch
	}

	return nil
}

//协程代理函数
func coroutine_proxy(ch chan int, current int, total int, f func(current int, total int)) {
	//catch panic error and free the resource
	defer func() {
		if err := recover(); err != nil {
			ch <- 1
			log.Fatal("Coroutine 'coroutine_proxy' Runtime Error: ", err)
		}
	}()
	f(current, total)
	ch <- 1
}
