package main

import (
	"flag"
	"runtime"
	"os"
	"fmt"
	"log"
	//"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"bytes"
	"io"
)

func Get(url string) (string,error) {

	// 超时时间：5秒
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "",err
		}
	}

	return result.String(),nil
}

func main() {
	var (
		//err error
		logFileName = flag.String("./log", "InstallMysql.log", "Log file name")
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "InstallMysql start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//write log
	log.Printf("Log start! File:%v \n", "InstallMysql.log")

	//r := gin.Default()
	//r.GET("/prometheus/api/v1/query",func(c *gin.Context){
	//	c.JSON(200,c.Query("query"))
	//	//if err != nil{
	//	//	log.Println(err.Error())
	//	//}
	//	//log.Println("AutoInStallMysql succeeded!")
	//})
	//r.Run(":80")

	//resp,err := http.Get("http://192.168.186.137/prometheus/api/v1/query?query=node_disk_io_now")
	//resp,err := http.Get("http://192.168.186.137/graph/api/datasources/proxy/1/api/v1/query_range?query=clamp_max(((avg%20by%20(mode)%20(%20(clamp_max(rate(node_cpu{instance%3D~%22localhost.localdomain%22%2Cmode!%3D%22idle%22}[5m])%2C1))%20or%20(clamp_max(irate(node_cpu{instance%3D~%22localhost.localdomain%22%2Cmode!%3D%22idle%22}[5m])%2C1))%20))*100%20or%20(avg_over_time(node_cpu_average{instance%3D~%22localhost.localdomain%22%2C%20mode!%3D%22total%22%2C%20mode!%3D%22idle%22}[5m])%20or%20avg_over_time(node_cpu_average{instance%3D~%22localhost.localdomain%22%2C%20mode!%3D%22total%22%2C%20mode!%3D%22idle%22}[5m])))%2C100)&start=1611169500&end=1611213000&step=300")
	//
	//if err!=nil{
	//	log.Println(err.Error())
	//}else{
	//	body,err := ioutil.ReadAll(resp.Body)
	//	if err!=nil {
	//		log.Println(err.Error())
	//	}
	//	fmt.Println(string(body))
	//}

	result,err:=Get("http://192.168.186.137/prometheus/api/v1/query?query=node_disk_io_now")
	if err!=nil {
		log.Println(err.Error())
	}else{
		fmt.Println(result)
	}

}
