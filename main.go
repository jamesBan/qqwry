package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
	"net"
	"errors"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	datFile := flag.String("qqwry", "./qqwry.dat", "纯真 IP 库的地址")
	port := flag.String("port", "2060", "HTTP 请求监听端口号")
	flag.Parse()

	IPData.FilePath = *datFile

	startTime := time.Now().UnixNano()
	res := IPData.InitIPData()

	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	endTime := time.Now().UnixNano()

	log.Printf("IP 库加载完成 共加载:%d 条 IP 记录, 所花时间:%.1f ms\n", IPData.IPNum, float64(endTime-startTime)/1000000)

	// 下面开始加载 http 相关的服务
	http.HandleFunc("/", findIP)

	log.Printf("开始监听网络端口:%s", *port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil); err != nil {
		log.Println(err)
	}
}

// findIP 查找 IP 地址的接口
func findIP(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(w, r)

	ip := r.Form.Get("ip")

	if ip == "" {
		res.ReturnError(http.StatusBadRequest,  "请填写 IP 地址")
		return
	}

	ips := strings.Split(ip, ",")
	err, ips := validateIp(ips)
	if err != nil {
		res.ReturnError(http.StatusBadRequest,  err.Error())
		return
	}

	if len(ips) <1 {
		res.ReturnError(http.StatusBadRequest, "请填写 IP 地址")
		return
	}

	qqWry := NewQQwry()

	rs := map[string]ResultQQwry{}
	if len(ips) > 0 {
		for _, v := range ips {
			rs[v] = qqWry.Find(v)
		}
	}

	res.ReturnSuccess(rs)
}

func validateIp(ips []string)(error, []string) {
	if len(ips) < 1 {
		return errors.New("ip 不能为空"), nil
	}

	filteredIps := []string{}
	for _, ip := range ips {
		if net.ParseIP(ip) != nil {
			filteredIps = append(filteredIps, ip)
		}
	}

	return nil, filteredIps
}