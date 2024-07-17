package main

import (
	"context"
	"flag"
	"fmt"
	"internal/cache"
	"log"
	"time"
)

type Args struct {
	cachType string
}

// arguments 의 기본값 지정 필요하면 사용
func NewArgs() Args {
	// return Args{cachType: cache.INMEMORYCACHE}
	return Args{}
}

func main() {

	args := NewArgs()
	flag.StringVar(&args.cachType, "t", args.cachType, "Cache Type (m : inmemory, r : redis ) [ default : dummy cache ]")
	flag.Parse()

	c := cache.NewCache(args.cachType)
	if c == nil {
		log.Fatalln("invalid cache type : ", args.cachType)
	}
	fmt.Println(c.Description())

	// 키값 저장 및 TTL 확인
	key := "k1"

	err := c.Set(context.TODO(), key, "in memory cache", 3*time.Second)
	if err != nil {
		log.Printf("'%s' %v\n", key, err)
	}

	v, err := c.Get(context.TODO(), key)
	if err != nil {
		log.Printf("'%s' %v\n", key, err)
	} else {
		fmt.Printf("%s key : %s , value : %+v\n", time.Now().Format("2006/01/02 15:04:05"), key, v)
	}

	ttl, err := c.GetTTL(context.TODO(), key)
	if err != nil {
		log.Printf("'%s' %v\n", key, err)
	} else {
		fmt.Printf("%s key : %s ttl : %v\n", time.Now().Format("2006/01/02 15:04:05"), key, ttl)
	}

	// 강제 key 값 만료 후 값 확인
	fmt.Println("5 second sleep ...")
	time.Sleep(5 * time.Second)

	v, err = c.Get(context.TODO(), key)
	if err != nil {
		log.Printf("'%s' %v\n", key, err)
	} else {
		fmt.Printf("%s key : %s , value : %+v\n", time.Now().Format("2006/01/02 15:04:05"), key, v)
	}

	ttl, err = c.GetTTL(context.TODO(), key)
	if err != nil {
		log.Printf("'%s' %v\n", key, err)
	} else {
		fmt.Printf("%s key : %s ttl : %v", time.Now().Format("2006/01/02 15:04:05"), key, ttl)
	}

}
