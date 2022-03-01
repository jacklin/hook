/*
 * @Descripttion:
 * @version:
 * @Author: linjincheng
 * @Date: 2021-12-03 17:39:27
 * @LastEditors: linjiancheng
 * @LastEditTime: 2021-12-03 18:00:25
 */
package hook

import (
	"testing"

	"github.com/jacklin/hook"
)

func TestSign(t *testing.T) {
	a := make(map[string]interface{})

	a["de"] = "11111"
	a["ab"] = "233aa"
	a["av"] = 100
	a["10"] = "232323"
	singv1 := &hook.SignV1{"key", "md5"}
	t.Log(singv1.GenerateSign(a))
}

func Benchmark_Add(b *testing.B) {
	a := make(map[string]interface{})

	a["de"] = "11111"
	a["ab"] = "233aa"
	a["av"] = 100
	a["10"] = "232323"
	singv1 := &hook.SignV1{"key", "md5"}
	for i := 0; i < b.N; i++ {
		singv1.GenerateSign(a)
	}
}
