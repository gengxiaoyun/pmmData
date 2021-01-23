package main

import (
	"flag"
	"runtime"
	"os"
	"fmt"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	simplejson "github.com/bitly/go-simplejson"
)

const(
	nodeCpuSecondsTotalRange = "http://192.168.186.137/prometheus/api/v1/query_range?query=" +
		"node_cpu_seconds_total{node_name=%22ubuntu%22}&" +
		"start=1611196430&end=1611196440&step=10"
	nodeCpuSecondsTotal = "http://192.168.186.137/prometheus/api/v1/query?query=" +
		"node_cpu_seconds_total{node_name=%22ubuntu%22}"
)

func HttpProcess(url string) (string,error) {
	req,err:=http.NewRequest("GET",url,nil)
	req.SetBasicAuth("admin","123456")
	if err!=nil{
		return "",err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp,err:=client.Do(req)

	if err!=nil{
		return "",err
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		return "",err
	}
	return string(body),nil
}

func JsonProcess(url string) error{
	str,err:=HttpProcess(url)
	if err!=nil{
		return err
	}
	res,err := simplejson.NewJson([]byte(str))
	if err != nil {
		return err
	}
	result,err := res.Get("data").Get("result").Array()
	if err != nil {
		return err
	}
	for _, row := range result {
		if each_map, ok := row.(map[string]interface{}); ok {
			if eachMapMetric, ok := each_map["metric"].(map[string]interface{}); ok {
				log.Printf("metric:{cpu=%s,mode=%s,node_name=%s}\n", eachMapMetric["cpu"],
					eachMapMetric["mode"], eachMapMetric["node_name"])
			}
			if eachMapValues, ok := each_map["values"].([]interface{}); ok {
				for i := 0; i < len(eachMapValues); i++ {
					log.Printf("values[%d]:{timestamp=%s,value=%s}\n",i,eachMapValues[i].([]interface{})[0],
						eachMapValues[i].([]interface{})[1])
				}
			}
			if eachMapValue, ok := each_map["value"].([]interface{}); ok {
				log.Printf("value:{timestamp=%s,value=%s}\n",eachMapValue[0],eachMapValue[1])
			}
		}
	}
	return nil
}

func main() {
	var (
		err error
		logFileName = flag.String("./log", "NodeExporter.log", "Log file name")
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//set logfile Stdout
	_, err = os.Stat(*logFileName)
	if err == nil {
		err = os.Remove(*logFileName)
		if err != nil {
			log.Println(err.Error())
		}
	}
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "logFile Open Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//write log
	log.Printf("Log start! File:%v \n", "NodeExporter.log")

	// query_range
	log.Println("query=node_cpu_seconds_total, start=1611196430, end=1611196440, step=10\n")
	err = JsonProcess(nodeCpuSecondsTotalRange)
	if err!=nil{
		log.Println(err.Error())
	}
	// query
	log.Println("query=node_cpu_seconds_total\n")
	err = JsonProcess(nodeCpuSecondsTotal)
	if err!=nil{
		log.Println(err.Error())
	}
}
