package hook

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

type SignV1 struct {
	Key      string `json:"key"`      //签名key
	SignType string `json:"signtype"` //签名类型 md5...
}

func (signv1 *SignV1) Sign(formatParams interface{}) string {
	//追加密钥
	sign := formatParams.(string) + signv1.Key
	//md5加密
	m := md5.New()
	m.Write([]byte(sign))
	sign = hex.EncodeToString(m.Sum(nil))
	return sign
}
func (signv1 *SignV1) formatParams(params interface{}) (string, error) {
	//解析为字节数组
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	//重组字符串
	var sign string
	newString := string(paramBytes)
	//为保证签名前特殊字符串没有被转码，这里解码一次
	newString = strings.Replace(newString, "\u003c", "<", -1)
	newString = strings.Replace(newString, "\u003e", ">", -1)

	//去除特殊标点
	newString = strings.Replace(newString, "\"", "", -1)
	newString = strings.Replace(newString, "{", "", -1)
	newString = strings.Replace(newString, "}", "", -1)

	paramArray := strings.Split(newString, ",")
	paramMap := make(map[interface{}]interface{})
	for _, v := range paramArray {
		detail := strings.SplitN(v, ":", 2)
		paramMap[detail[0]] = detail[1]
	}
	lm := keySort(paramMap, "string")
	var arr []string
	for _, v := range lm {
		arr = append(arr, v.k+":"+v.v)
	}
	for _, v := range arr {
		detail := strings.SplitN(v, ":", 2)
		//排除sign和sign_type
		if detail[0] != "sign" && detail[0] != "sign_type" {
			if sign == "" {
				sign = detail[0] + "=" + detail[1]
			} else {
				sign += "&" + detail[0] + "=" + detail[1]
			}
		}
	}

	return sign, nil
}

func (signv1 *SignV1) GenerateSign(params interface{}) string {
	if fp, err := signv1.formatParams(params); err != nil {
		return ""
	} else {
		sign := signv1.Sign(fp)
		return sign
	}
}

type Lmap struct {
	k string
	v string
}

func keySort(data map[interface{}]interface{}, sortType interface{}) []Lmap {
	var t []interface{}
	for k, _ := range data {
		t = append(t, k)
	}
	var st []string
	var it []int

	if reflect.TypeOf(sortType).String() == "int" {
		for _, v := range t {
			it = append(it, v.(int))
		}
		sort.Ints(it)
	} else if reflect.TypeOf(sortType).String() == "string" {
		for _, v := range t {
			st = append(st, v.(string))
		}
		sort.Strings(st)
	}
	var lm []Lmap
	if it != nil {
		for _, k := range it {
			lm = append(lm, Lmap{string(k), data[k].(string)})
		}
	} else {
		for _, k := range st {
			lm = append(lm, Lmap{k, data[k].(string)})
		}
	}
	return lm
}
