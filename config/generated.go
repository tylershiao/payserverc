package config

import (
	"io"
	"io/ioutil"
	"os"
	"server/util"
	"strconv"
	"strings"
)

var (
	luas = []string{"cn", "en"} // 支持的语言版本
)

func BuildErrCode() {
	f, err := ioutil.ReadFile(".." + string(os.PathSeparator) + "pb" + string(os.PathSeparator) + "proto" + string(os.PathSeparator) + "common.proto")
	if err != nil {
		panic("read fail")
	}
	proto := string(f)
	eIndexS := strings.Index(proto, "ErrNo")
	tmp := proto[eIndexS:]
	eIndexS2 := strings.Index(tmp, "{")
	tmp = tmp[eIndexS2:]
	eIndexE := strings.Index(tmp, "}")
	target := tmp[:eIndexE+1]
	lines := strings.Split(target, "\n")

	var filename = "errcode.toml"
	if util.CheckFileIsExist(filename) { //如果文件存在
		os.Remove(filename)
	}
	file, err := os.Create(filename) //创建文件
	if err != nil {
		return
	}
	// file, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
	errNums := []string{}
	errDes := make([][]string, len(luas))
	for _, v := range lines {
		indexS := strings.Index(v, "=")
		indexE := strings.Index(v, ";")
		dIndexS := strings.Index(v, "//")
		if indexS < 0 || indexE < 0 || dIndexS < 0 {
			continue
		}
		numS := v[indexS+1 : indexE]
		num := strings.Trim(numS, " ")
		if tmp, err := strconv.Atoi(num); err != nil {
			return
		} else if tmp == 0 || tmp == 200 {
			continue
		}
		errNums = append(errNums, num)

		if !strings.Contains(v[dIndexS+2:], ";") {
			panic("config err")
		}

		allDes := strings.Split(v[dIndexS+2:], ";")
		for i := range allDes {
			if i >= len(errDes) {
				panic("config err")
			}
			errDes[i] = append(errDes[i], allDes[i])
		}
	}
	for i, v := range luas {
		if _, err := io.WriteString(file, "["+v+"]"+"\n"); err != nil {
			panic("config err")
		}
		for j, v := range errDes[i] {
			des1 := strings.Trim(v, " ")
			des2 := strings.Trim(des1, string(byte(13)))
			writeS := errNums[j] + " =" + " " + `"` + des2 + `"` + "\n"
			if j == len(errDes[i])-1 {
				writeS += "\n"
			}
			if _, err := io.WriteString(file, writeS); err != nil {
				panic("config err")
			}
		}
	}
}
