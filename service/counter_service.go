package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// CounterHandler 计数器接口
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理计数器请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	res := &JsonResult{}

	if r.Method == http.MethodGet {
		LogStep("处理GET请求，获取当前计数器", nil)
		counter, err := getCurrentCounter()
		LogDBResult("查询", "counters", counter, err)
		if err != nil {
			LogError("获取计数器失败", err)
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			LogStep("获取计数器成功", map[string]interface{}{"count": counter.Count})
			res.Data = counter.Count
		}
	} else if r.Method == http.MethodPost {
		LogStep("处理POST请求，修改计数器", nil)
		count, err := modifyCounter(r)
		if err != nil {
			LogError("修改计数器失败", err)
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			LogStep("修改计数器成功", map[string]interface{}{"count": count})
			res.Data = count
		}
	} else {
		LogError("请求方法不支持", fmt.Errorf("期望GET或POST方法，实际为%s", r.Method))
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		LogError("序列化响应失败", err)
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
	LogResponse(res, nil)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	LogStep("开始解析action参数", nil)
	action, err := getAction(r)
	if err != nil {
		LogError("获取action参数失败", err)
		return 0, err
	}
	LogStep("获取action参数成功", map[string]string{"action": action})

	var count int32
	if action == "inc" {
		LogStep("执行自增操作", nil)
		count, err = upsertCounter(r)
		if err != nil {
			LogError("自增操作失败", err)
			return 0, err
		}
		LogStep("自增操作成功", map[string]interface{}{"count": count})
	} else if action == "clear" {
		LogStep("执行清零操作", nil)
		err = clearCounter()
		if err != nil {
			LogError("清零操作失败", err)
			return 0, err
		}
		count = 0
		LogStep("清零操作成功", map[string]interface{}{"count": count})
	} else {
		LogError("action参数错误", fmt.Errorf("不支持的action: %s", action))
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	LogStep("开始查询当前计数器", nil)
	currentCounter, err := getCurrentCounter()
	LogDBResult("查询", "counters", currentCounter, err)

	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		LogError("查询当前计数器失败", err)
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		LogStep("计数器不存在，创建新计数器", nil)
		count = 1
		createdAt = time.Now()
	} else {
		LogStep("计数器存在，执行自增", map[string]interface{}{"currentCount": currentCounter.Count})
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	LogStep("准备更新计数器", map[string]interface{}{"id": counter.Id, "count": counter.Count})
	LogDBOperation("更新", "counters", counter)
	err = dao.Imp.UpsertCounter(counter)
	if err != nil {
		LogError("更新计数器失败", err)
		return 0, err
	}
	LogDBResult("更新", "counters", counter, nil)
	LogStep("计数器更新成功", map[string]interface{}{"count": counter.Count})
	return counter.Count, nil
}

func clearCounter() error {
	LogStep("开始清零计数器", map[string]interface{}{"counterId": 1})
	LogDBOperation("清零", "counters", map[string]interface{}{"id": 1})
	err := dao.Imp.ClearCounter(1)
	LogDBResult("清零", "counters", nil, err)
	if err != nil {
		LogError("清零计数器失败", err)
	} else {
		LogStep("计数器清零成功", nil)
	}
	return err
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	LogStep("查询当前计数器", map[string]interface{}{"counterId": 1})
	LogDBOperation("查询", "counters", map[string]interface{}{"id": 1})
	counter, err := dao.Imp.GetCounter(1)
	LogDBResult("查询", "counters", counter, err)
	return counter, err
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
	LogStep("开始解析请求体", nil)
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		LogError("解析请求体失败", err)
		return "", err
	}
	defer r.Body.Close()

	LogStep("请求体解析成功", body)
	action, ok := body["action"]
	if !ok {
		LogError("缺少action参数", fmt.Errorf("请求体中未找到action字段"))
		return "", fmt.Errorf("缺少 action 参数")
	}

	LogStep("获取action参数成功", map[string]string{"action": action.(string)})
	return action.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
