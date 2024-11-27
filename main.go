package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/miekg/dns"
)

func checkPort(ip string, port string) bool {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("baidu.com."), dns.TypeA) // set the question to ask
	m.RecursionDesired = true

	msg, td, err := c.Exchange(m, ip+":"+port)
	fmt.Printf("msg: %s, td: %f\n", msg.String(), td.Seconds())
	return (len(msg.String()) > 0 || td.Seconds() > 0) && err == nil
}

func worker(records <-chan []string, results chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for record := range records {
		ip := record[0]             // 假设IP在第一列
		landuiServerId := record[1] // 假设服务器ID在第二列
		var status string
		if checkPort(ip, "53") {
			fmt.Println(ip, "is open")
			status = "open"
		} else {
			fmt.Println(ip, "is closed")
			status = "closed"
		}
		results <- []string{ip, landuiServerId, status}
	}
}

func main() {
	inputFilePath := flag.String("input", "ip_list.csv", "Path to the input CSV file")
	outputFilePath := flag.String("output", "output.csv", "Path to the output CSV file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println(os.Args[0], "-input=input.csv -output=output.csv")
	}
	flag.Parse()

	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)

	// 读取CSV文件的所有记录
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 写入CSV文件的表头
	writer.Write([]string{"IP", "Server ID", "DNS status"})

	// 创建一个channel来传递记录
	recordChan := make(chan []string, len(records))
	resultsChan := make(chan []string, len(records))

	// 创建一个WaitGroup来等待所有worker goroutine完成
	var wg sync.WaitGroup

	// 启动固定数量的worker goroutine
	numWorkers := 50
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(recordChan, resultsChan, &wg)
	}

	// 将记录发送到channel中
	go func() {
		for _, record := range records[1:] { // 跳过表头
			recordChan <- record
		}
		close(recordChan)
	}()

	// 启动一个goroutine来写入CSV文件
	go func() {
		for result := range resultsChan {
			writer.Write(result)
		}
	}()

	// 等待所有worker goroutine完成
	wg.Wait()
	close(resultsChan)
}
