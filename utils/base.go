package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	Y_M = "2006-01"
	Y_M_D = "2006-01-02"
	Y_M_D_2 = "2006年01月02日"
	Y_M_D_H_I_S = "2006-01-02 15:04:05"
	Y_M_D_H_I_S_2 = "2006年01月02日 15:04:05"
	H_I_S = "15:04:05"
)

func Format(str interface{}, layout string) string {
	var date time.Time
	var err error
	//判断变量类型
	switch str.(type) {
	case time.Time:
		date = str.(time.Time)
	case string:
		//如果是字符串则转换成 标准日期时间格式
		fmt.Println(str)
		date, err = time.Parse(layout, str.(string))
		if err != nil {
			return ""
		}
	}
	return date.Format(layout)
}

//当前日期时间
func Now() string {
	return Format(time.Now(), Y_M_D_H_I_S)
}
//当前日期
func Date() string {
	return Format(time.Now(), Y_M_D)
}
//当前时间
func Time() string {
	return Format(time.Now(), H_I_S)
}

//字符串base64加密
func Base64E(urlstring string) string {
	str := []byte(urlstring)
	data := base64.StdEncoding.EncodeToString(str)
	return data
}

//字符串base64解密
func Base64D(urlxxstring string) string {
	data, err := base64.StdEncoding.DecodeString(urlxxstring)
	if err != nil {
		return ""
	}
	s := fmt.Sprintf("%q", data)
	s = strings.Replace(s, "\"", "", -1)
	return s
}

//url转义
func UrlE(s string) string {
	return url.QueryEscape(s)
}

//url解义
func UrlD(s string) string {
	s, e := url.QueryUnescape(s)
	if e != nil {
		return e.Error()
	} else {
		return s
	}
}

//得到系统时间
func GetTime() time.Time {
	timezone := float64(0)
	timezone, _ = strconv.ParseFloat("8", 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add))
}

/*"2006-01-02 15:04:05"*/
//得到今天日期字符串
func GetTodayString() string {
	timezone := float64(0)
	timezone, _ = strconv.ParseFloat("8", 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add)).Format("20060102")
}

//得到时间字符串
func GetTimeString() string {
	timezone := float64(0)
	timezone, _ = strconv.ParseFloat("8", 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add)).Format("20060102150405")
}

type Map map[string]interface{}

func (newMap Map) Merge(oldMap Map) {
	for k, v := range oldMap {
		newMap[k] = v
	}
}

/**
 * 根据path读取文件中的内容，返回字符串
 * 建议使用绝对路径，例如："./schema/search/appoint.json"
 */
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func ReadJson(path string) Map {
	jsonStr := ReadFile(path)
	ret := Map{}
	err := json.Unmarshal([]byte(jsonStr), &ret)
	if err != nil {
		panic("文件[" + path + "]的内容不是json格式")
	}
	return ret
}

func JsonToMap(jsonStr string) Map {
	var mapResult Map
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		fmt.Println("JsonToMap err: ", err)
	}
	return mapResult
}

/**
arrayStr为原始字符串，elementType为输出数组元素类型
*/
func StringToArray(arrayStr, elementType string) interface{} {
	var resultArray interface{}
	var err error
	switch elementType {
	case reflect.String.String():
		var tempArray []string
		err = json.Unmarshal([]byte(arrayStr), &tempArray)
		resultArray = tempArray
	case reflect.Int.String():
		var tempArray []int
		err = json.Unmarshal([]byte(arrayStr), &tempArray)
		resultArray = tempArray
	case reflect.Map.String():
		var tempArray []map[string]interface{}
		err = json.Unmarshal([]byte(arrayStr), &tempArray)
		resultArray = tempArray
	default:
		var tempArray []interface{}
		err = json.Unmarshal([]byte(arrayStr), &tempArray)
		resultArray = tempArray
	}
	if err != nil {
		fmt.Println("string to array failure: ", err.Error())
		return nil
	}
	return resultArray
}

func StructToMap(source interface{}, withReflect bool) map[string]interface{} {
	target := make(Map)
	if withReflect {
		relVal := reflect.ValueOf(source)
		relType := reflect.TypeOf(source)
		if relType.Kind() != reflect.Struct {
			fmtString := fmt.Sprintf("待转换类型(%s)不为结构体类型, 参数错误.", relType.String())
			fmt.Println(fmtString)
			return target
		}
		for i := 0; i < relVal.NumField(); i++ {
			if relType.Field(i).Name == "BaseModel" {
				target.Merge(StructToMap(relVal.Field(i).Interface(), true))
			} else {
				key := relType.Field(i).Name
				target[key] = relVal.Field(i).Interface()
			}
		}
		return target
	}
	jsonStr, _ := json.Marshal(source)
	_ = json.Unmarshal(jsonStr, &target)
	return target
}

/**
 * 对象转换为int
 * 支持类型：int,float64,string,bool(true:1;false:0)
 * 其他类型报错
 */
func ToInt(obj interface{}) int {
	switch obj.(type) {
	case int:
		return obj.(int)
	case float64:
		return int(obj.(float64))
	case string:
		ret, _ := strconv.Atoi(obj.(string))
		return ret
	case bool:
		if obj.(bool) {
			return 1
		} else {
			return 0
		}
	default:
		panic("ToInt出错")
	}
}

/**
 * 对象转换为string
 * 支持类型：int,float64,string,bool(true:"1";false:"0")
 * 其他类型报错
 */
func ToString(obj interface{}) string {
	switch obj.(type) {
	case int:
		return strconv.Itoa(obj.(int))
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case string:
		return obj.(string)
	case bool:
		if obj.(bool) {
			return "1"
		} else {
			return "0"
		}
	default:
		panic("ToString出错")
	}
}

/**
 * 对象转换为bool
 * 支持类型：int,float64,string,bool
 * 其他类型报错
 */
func ToBool(obj interface{}) bool {
	switch obj.(type) {
	case int:
		if obj.(int) == 0 {
			return false
		} else {
			return true
		}
	case float64:
		if obj.(float64) == 0 {
			return false
		} else {
			return true
		}
	case string:
		trues := map[string]int{"true": 1, "是": 1, "1": 1, "真": 1}
		if _, ok := trues[strings.ToLower(obj.(string))]; ok {
			return true
		} else {
			return false
		}
	case bool:
		return obj.(bool)
	default:
		panic("ToBool出错")
	}
}

//判断一个数据是否为空，支持int, float, string, slice, array, map的判断
func Empty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if reflect.ValueOf(value).Len() == 0 {
			return true
		} else {
			return false
		}
	}
	return false
}

//判断某一个值是否在列表(支持 slice, array, map)中
func InList(needle interface{}, haystack interface{}) bool {
	//interface{}和interface{}可以进行比较，但是interface{}不可进行遍历
	hayValue := reflect.ValueOf(haystack)
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice, reflect.Array:
		//slice, array类型
		for i := 0; i < hayValue.Len(); i++ {
			if hayValue.Index(i).Interface() == needle {
				return true
			}
		}
	case reflect.Map:
		//map类型
		var keys []reflect.Value = hayValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if hayValue.MapIndex(keys[i]).Interface() == needle {
				return true
			}
		}
	default:
		return false
	}
	return false
}
