package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/load"
	"os/exec"
	"time"
)


func init() {
	flag.StringVar(&testfile, "testfile", "/usr/share/sysbench/oltp_read_write.lua", "是否使用默认测试用例")
	flag.IntVar(&loadBase, "loadBase", 4, "测试基础负载，uptime小于改值时开始测试")
	flag.StringVar(&host, "host", "127.0.0.1", "数据库地址")
	flag.StringVar(&user, "user", "root", "数据库用户")
	flag.StringVar(&password, "password", "123456", "数据库密码")

}
var user string
var password string
var testfile string
var host string
var loadBase int

const  sysbanchCmd =
	"/usr/bin/sysbench %s  --max-requests=10000 --num-threads=%d --mysql-host=%s --mysql-port=3306 --mysql-user=%s " +
	"--mysql-db=sbtest --mysql-password=%s run " +
	"|grep transactions |awk '{print $3}'"

var  threadNums []int

func main(){
	flag.Parse()
	fmt.Println("sysbench mysql")
	threadNums = []int{1,10,50,100,200,500,1000,2000}
	for i:=0 ;i< len(threadNums);{
		thread := threadNums[i]
		time.Sleep(time.Second * 1)
		load,err := load.Avg()
		if err != nil {
			fmt.Print("...")
			continue
		}
		if load.Load1 > float64(loadBase){
			fmt.Print("...")
			continue
		}
		fmt.Println("...")
		cmdStr := fmt.Sprintf(sysbanchCmd,testfile,thread,host,user,password)
		cmd := exec.Command("sh", "-c",cmdStr)
		result,err := cmd.CombinedOutput()
		if err!= nil{
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("thread: %d tps: %s",thread,string(result))
		i++
	}
}

