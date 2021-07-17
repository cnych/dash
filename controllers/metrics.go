package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/cnych/dash/models"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func GetMetrics(c *gin.Context) {
	// POST 传递的 JSON 数据
	var metricsQuery models.MetricsQuery
	if err := c.ShouldBindJSON(&metricsQuery); err != nil {
		klog.V(2).ErrorS(err, "bind models.MetricsQuery to json failed", "controller", "LoadMetrics")
		writeOK(c, gin.H{})
		return
	}

	// todo，应该从数据库中获取 Prometheus 的服务
	// 先判断下 Prometheus 服务是否可用（配置一个超时参数）
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	readyReq, err := http.NewRequest("GET", "http://192.168.31.30:32165/-/ready", nil)
	if err != nil {
		klog.V(2).ErrorS(err, "check prometheus service ready request failed", "controller", "LoadMetrics")
		writeOK(c, gin.H{})
		return
	}
	readyResp, err := http.DefaultClient.Do(readyReq.WithContext(ctx))
	if err != nil {
		klog.V(2).ErrorS(err, "check prometheus service ready failed", "controller", "LoadMetrics")
		writeOK(c, gin.H{})
		return
	}
	// 如果还没有ready，则直接返回前端空数据
	if readyResp.StatusCode != http.StatusOK {
		writeOK(c, gin.H{})
		return
	}

	step := 60
	end := time.Now().Unix()
	start := end - 3600

	// tracker
	tracker := models.NewPrometheusTracker()
	wg := sync.WaitGroup{}

	e := reflect.ValueOf(&metricsQuery).Elem()
	for i := 0; i < e.NumField(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			fName := e.Type().Field(i).Name
			fValue := e.Field(i).Interface().(*models.MetricsCategory)
			fTag := e.Type().Field(i).Tag
			if fValue == nil {
				return
			}
			klog.V(3).InfoS("start request prometheus data", "filed", fName, "controller", "LoadMetrics")
			// 请求 Prometheus 查询
			prometheusQueries := fValue.GenerateQuery()
			if prometheusQueries == nil {
				klog.V(2).InfoS("no promql", "field", fName, "controller", "LoadMetrics")
				return
			}
			promql := url.QueryEscape(prometheusQueries.GetValueByField(fName))
			klog.V(2).InfoS("prometheus query ql", "field", fName, "promql", promql)
			// todo，从数据库中获取prometheus地址
			resp, err := http.Get(fmt.Sprintf("http://192.168.31.30:32165/api/v1/query_range?query=%s&start=%d&end=%d&step=%d", promql, start, end, step))
			if err != nil {
				klog.V(2).ErrorS(err, "request metrics data failed", "controller", "LoadMetrics")
				return
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				klog.V(2).ErrorS(err, "read response body failed", "controller", "LoadMetrics")
				return
			}
			var data models.PrometheusQueryResp
			if err := json.Unmarshal(body, &data); err != nil {
				klog.V(2).ErrorS(err, "unmarshal response body to models.PrometheusQueryResp failed",
					"controller", "LoadMetrics")
				return
			}
			// 配置当前查询的数据结果
			tag := fTag.Get("json")
			tracker.Set(tag[:strings.Index(tag, ",omitempty")], &data)
		}(i)
	}
	// 等待所有查询完成
	wg.Wait()

	// 返回最后拼接的数据
	writeOK(c, gin.H{
		"metrics": tracker.Metrics,
	})
}

