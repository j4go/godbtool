
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/akkuman/parseConfig"
	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

var (
	elasticUrl  string
	docType     string
	elasticUser string
	elasticPwd  string
	// logtoolUser    string
	// logtoolPwd     string
	requestTimeOut float64
)

func getDateStr() string {
	loc, _ := time.LoadLocation("Local") //服务器设置的时区
	now := time.Now().In(loc)
	t := "2006.01.02"
	dateStr := now.Format(t)
	return dateStr
}

func getLogFileName() string {
	loc, _ := time.LoadLocation("Local") //服务器设置的时区
	now := time.Now().In(loc)
	t := "20060102_1504"
	dateStr := now.Format(t)
	return "gin_" + dateStr + ".log"
}

func logger(logContent interface{}) {
	fmt.Fprintln(gin.DefaultWriter, logContent)
	return
}

func elasticRequestHelper(url string, jsonStr string) string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	auth := []string{elasticUser, elasticPwd}
	timeOut := 5000 * time.Millisecond
	options := &grequests.RequestOptions{
		Headers:        headers,
		JSON:           jsonStr,
		Auth:           auth,
		RequestTimeout: timeOut,
	}
	resp, err := grequests.Post(url, options)
	if err != nil {
		return err.Error()
	} else {
		return resp.String()
	}
}

func logHandler(c *gin.Context) {

	// get router as path
	path := c.Param("path")
	logger("path==" + path)

	// log type elastic or more
	// logType := c.DefaultQuery("log_type", "elastic")
	// elUrl := c.DefaultQuery("el_url", "")

	// Parse JSON
	var json_struct map[string]interface{}

	if c.BindJSON(&json_struct) == nil {
		msg := make(map[string]interface{})
		msg["msg"] = json_struct // body
		// logger(reflect.TypeOf(json_struct))

		// 预留字段处理 剪切出来放到上一级
		m := msg["msg"].(map[string]interface{})

		if addr, ok := m["$addr"]; ok {
			msg["$addr"] = addr
			delete(m, "$addr")
		}
		if geo, ok := m["$geo"]; ok {
			msg["$geo"] = geo
			delete(m, "$geo")
		}
		if dev_type, ok := m["$dev_type"]; ok {
			msg["$dev_type"] = dev_type
			delete(m, "$dev_type")
		}

		// headerMap := c.Request.Header
		logger(c.Request.Header)

		ct := c.GetHeader("X-Rf-Data-Type")
		if ct != "json" {
			c.JSON(400, gin.H{"status": "bad request data type"})
			return
		}

		// 以下三个是nginx带过来的
		x_forward_for := c.GetHeader("x_forward_for")
		if x_forward_for != "" {
			msg["$x_forward_for"] = x_forward_for
		}
		x_rf_path := c.GetHeader("x_rf_path")
		if x_rf_path != "" {
			msg["$path"] = x_rf_path
		}
		x_rf_date := c.GetHeader("x_rf_date") // 时间戳字符串
		if x_rf_date != "" {
			x_rf_date_int64, _ := strconv.ParseInt(x_rf_date, 10, 64)
			msg["$log_timestamp"] = time.Unix(x_rf_date_int64, 0) // x_rf_date:=time.Now().Unix() 单位秒
		}

		// 系统接收到的时间
		loc, _ := time.LoadLocation("Local") //服务器设置的时区
		now := time.Now().In(loc)
		logger(reflect.TypeOf(now))
		msg["@timestamp"] = now // 系统接收时间 time.Time 可能需要转化

		// 需要根据index和template限制验证一些字段的格式比如date等

		logger(msg)

		mjson, err := json.Marshal(msg)
		if err != nil {
			logger(err)
		}
		jsonStr := string(mjson)

		logger(jsonStr)

		var res string

		// 根据path的值选择不同的操作 先写死了  后面更多的值再选择
		if path == "lgf-test" || path == "dev-test" || path == "monitor" || path == "voltmeter-prod" || path == "voltmeter-test" {
			indexName := "iot-" + path + "-" + getDateStr() + "-test"
			logger(indexName)
			url := elasticUrl + "/" + indexName + docType
			logger(url)
			res = elasticRequestHelper(url, jsonStr)
			logger("elastic result:" + res)
			resMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(res), &resMap)
			if err == nil {
				if _, ok := resMap["_shards"]; ok {
					logger(resMap["_shards"])
					c.JSON(http.StatusOK, gin.H{"status": "ok"})
				} else {
					// json结果格式不符合
					logger(resMap)
					c.JSON(http.StatusOK, gin.H{"status": "json format unexpected"})
				}
			} else {
				// 出错 超时等
				logger(err)
				c.JSON(http.StatusOK, gin.H{"status": "unexpected result"})
			}
		} else if path == "xxx" {
			c.JSON(200, gin.H{"status": "ok"})
		} else {
			// path 不合规
			c.JSON(400, gin.H{"status": "bad request, try the right path"})
		}
	} else {
		// 缺少request body
		logger("缺少request body")
		c.JSON(400, gin.H{"status": "bad request data format"})
	}
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	logtoolUser: logtoolPwd,
	// }))

	r.PUT("/:path", logHandler)

	return r
}

func main() {

	// 从命令行读取参数 如环境变量 从文件读取配置
	var m string
	var c string
	var l string
	var config parseConfig.Config
	flag.StringVar(&m, "m", "release", "set gin env debug/release")
	flag.StringVar(&c, "c", "config_release.json", "set config file")
	flag.StringVar(&l, "l", "logs/", "set logs directory")
	flag.Parse()

	config = parseConfig.New(c)

	// gin.DisableConsoleColor()
	f, _ := os.Create(l + getLogFileName())
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	logger("\n===")
	logger("env: " + m)
	logger("config: " + c)
	logger("log directory: " + l)
	logger("===\n")

	// if m == "dev" {
	// 	config = parseConfig.New("conf/config_dev.json")
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	config = parseConfig.New("conf/config_release.json")
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// 读取配置文件 https://github.com/akkuman/parseConfig
	elasticUrl = config.Get("elasticUrl").(string)
	docType = config.Get("docType").(string)
	elasticUser = config.Get("elasticUser").(string)
	elasticPwd = config.Get("elasticPwd").(string)
	// logtoolUser = config.Get("logtoolUser").(string)
	// logtoolPwd = config.Get("logtoolPwd").(string)

	r := setupRouter()
	r.Run(config.Get("port").(string))
}
