//v1 只支持读kv格式
package yaml2map

import (
	"io/ioutil"
	"log"
	"regexp"
)

type configMap map[string]map[string]string
type kvMap map[string]string

var reComm = regexp.MustCompile("(^#.*)|(\\s*#.*)") //使用regexp库的Replace函数，去除注释
var reConf = regexp.MustCompile("-{3}")             //使用regexp库的Split函数，用---分割业务配置
var ren = regexp.MustCompile("\\n")                 //使用regexp库的Replace函数, 去除换行符
var reKvStr = regexp.MustCompile("\\s{2}")          //使用regexp库的Split函数, 分割kv string为kv slice
var reKv = regexp.MustCompile(":\\s")               //使用regexp库的Split函数, 分割kv 为slice, slice[0]为key slice[1]为value
var reSer = regexp.MustCompile(":\\s{2}")           //使用regexp库的Split函数, 分割业务配置为业务名,和kv string
var revs = regexp.MustCompile("\"")                 //使用regexp库的Replace函数, 去value中的"
// var reIsKvMap = regexp.MustCompile("^\\s{2}")       //使用regexp库的MatchString函数, 判断文件是configmap还是kvmap

func readFile(file string) string {
	byteFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("read config file %s error. err = %v\n", file, err)
	}
	strFile := string(byteFile)
	var noCommStrFile string
	// 配置文件存在注释，则替换为""
	if reComm.MatchString(strFile) {
		noCommStrFile = reComm.ReplaceAllString(strFile, "")
	} else {
		noCommStrFile = strFile
	}
	return noCommStrFile
}
func ReadConfigMap(file string) configMap {
	noCommStrFile := readFile(file)
	var confSlice []string
	// 配置文件存在多个服务配置，默认用“---”划分范围，按---分割
	if reConf.MatchString(noCommStrFile) {
		confSlice = reConf.Split(noCommStrFile, -1)
	} else {
		//配置文件只存在单个服务配置，不包含“---”，或没用“---”直接转成长度为1的切片
		confSlice = []string{noCommStrFile}
	}
	//根据切片长度创建返回的configmap
	confmap := make(configMap, len(confSlice))
	// 轮询配置文件切片元素，获取各个配置的配置名和kvmap
	for _, oconf := range confSlice {
		conf := ren.ReplaceAllString(oconf, "")
		confSlice := reSer.Split(conf, -1)
		kvSlices := reKvStr.Split(confSlice[1], -1)
		kvmap := make(kvMap, len(kvSlices))
		for _, kvs := range kvSlices {
			kv := reKv.Split(kvs, -1)
			kvmap[kv[0]] = revs.ReplaceAllString(kv[1], "")
		}
		confmap[confSlice[0]] = kvmap
	}
	return confmap
}
