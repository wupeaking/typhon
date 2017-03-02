package main

import (
	"testing"
	"regexp"
)

func TestRe(t *testing.T) {
	re := regexp.MustCompile(`/v1/networks/(?P<id>\d*)/(?P<name>\w*)`)
	suc := re.MatchString("/v1/networks/31313/dahdahdaiuhda")
	t.Log("是否匹配成功: ", suc)
	t.Log("分组名:", re.SubexpNames(), "数量: ", len(re.SubexpNames()))
	mat := re.FindStringSubmatch("/v1/networks/31313/dahdahdaiuhda")
	//result := re.FindSubmatch([]byte("/v1/networks/31313/dahdahdaiuhda"))
	for _, v := range mat{
		t.Log("match:", v)
	}
	//t.Log("len:", len(mat[2]))
	//t.Log("子匹配: ", re.("/v1/networks/121212/dahdahdaiuhda", 2))
	t.Log("分组数量: ", re.NumSubexp())

	newre:= regexp.MustCompile(`^/$`)
	ret := newre.FindStringSubmatch("/hfeuihfu")

	t.Log("匹配成功", ret, newre.MatchString("/feefefefefe"))
}