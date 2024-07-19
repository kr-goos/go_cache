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
	cacheType string
	addr      string
	password  string
	db        int
}

func NewDefaultArgs() Args {
	return Args{
		cacheType: cache.DUMMYCACHE,
		addr:      "localhost:6379",
	}
}

func main() {

	args := NewDefaultArgs()
	flag.StringVar(&args.cacheType, "t", args.cacheType, "Cache Type (memory : inmemory, redis : redis ) [ default : dummy cache ]")
	flag.StringVar(&args.addr, "a", args.addr, "server address (format: 8.8.8.8:80 )")
	flag.StringVar(&args.password, "p", args.password, "server password")
	flag.IntVar(&args.db, "d", args.db, "Database to be selected after connecting to the server")
	flag.Parse()

	c, err := cache.NewCache(args.cacheType, args.addr, args.password, args.db)
	if err != nil {
		log.Fatalln("invalid cache type : ", args.cacheType)
	}
	defer c.Close()

	fmt.Println(c.Description())

	// 키값 저장 및 TTL 확인
	key := "k1"

	err = c.Set(context.TODO(), key, "in memory cache", 3*time.Second)
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
