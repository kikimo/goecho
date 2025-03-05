// client.go
package main

import (
	"log"
	"math"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"sort"
	"strconv"
	"time"
)

// Args holds arguments passed to the RPC method
type Args struct {
	A, B int
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	args := &Args{A: 7, B: 8}
	var reply int

	delay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Duration(delay) * time.Microsecond)
	defer ticker.Stop()

	var latencies []time.Duration

	for i := 0; i < 1000; i++ { // 发送 1000 次请求
		start := time.Now()

		err = client.Call("Arith.Add", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}

		elapsed := time.Since(start)
		latencies = append(latencies, elapsed)

		<-ticker.C
	}

	// 计算统计数据
	mean := calculateMean(latencies)
	p50 := calculatePercentile(latencies, 50)
	p90 := calculatePercentile(latencies, 90)

	log.Printf("Mean: %s, P50: %s, P90: %s\n", mean, p50, p90)
}

// 计算均值
func calculateMean(latencies []time.Duration) time.Duration {
	var sum time.Duration
	for _, latency := range latencies {
		sum += latency
	}
	return time.Duration(int64(sum) / int64(len(latencies)))
}

// 计算百分位数
func calculatePercentile(latencies []time.Duration, percentile int) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	index := int(math.Ceil(float64(percentile)/100*float64(len(latencies))) - 1)
	return latencies[index]
}
