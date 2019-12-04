/**
 * @Author: DollarKillerX
 * @Description: context.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 上午9:24 2019/12/3
 */
package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	//timeout, _ := context.WithTimeout(context.TODO(), time.Second)
	//go prints(timeout)

	go prints1()
	time.Sleep(time.Second * 10)
}

func prints(ctx context.Context) {
	once := sync.Once{}
	for {
		select {
		case <-ctx.Done():
			log.Println("结束")
			return
		default:
			once.Do(func() {
				log.Println("ok")
				time.Sleep(time.Millisecond * 200)
			})
		}
	}
}

func prints1() {
	select {
	case <-time.After(time.Millisecond * 200):
		log.Println("time.after")
		return
	default:
		log.Println("ok")
		time.Sleep(time.Second * 2)
	}
}

//func prints2(ctx context.Context) error {
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				log.Println("dome")
//				return errors.New("time out")
//			default:
//			}
//		}
//	}()
//
//	log.Println("ac")
//	return nil
//}
