package smells

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func IsSharedDataManagement(urls []string) bool {

	//search for protocol
	re := regexp.MustCompile("jdbc:(.+?):")
	var ds []string
	for _, url := range urls {
		ds = re.FindStringSubmatch(url)
		if ds != nil {
			break
		}
	}

	//Filter repeating elements with a map characteristics of unique primary key
	tempMap := map[string]byte{}

	switch ds[1] {
	case "h2":
		for _, url := range urls {
			l := len(tempMap)
			tempMap[url] = 0
			if len(tempMap) == l {
				return true
			}
		}
	case "mysql":
		for _, url := range urls {
			re := regexp.MustCompile("jdbc:mysql://(.+?)/")
			match := re.FindStringSubmatch(url)
			l := len(tempMap)
			tempMap[match[1]] = 0
			if len(tempMap) == l {
				return true
			}
		}
	}

	return false
}

func ScanSDM(services []os.FileInfo, path string) ([]string, []string, []string) {
	log.Println("task -> scan shared data management")

	var dataSourceDev []string
	var dataSourceProd []string
	var dataSourceDocker []string
	for _, f := range services {
		log.Println("scanning", f.Name())

		data, _ := ioutil.ReadFile(path + "/" + f.Name() + "/src/main/resources/config/application-dev.yml")
		m := make(map[interface{}]interface{})
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource := m["spring"].(map[interface{}]interface{})["datasource"].(map[interface{}]interface{})["url"].(string)
		dataSource = strings.Replace(dataSource, "file:./", "file:"+path+"/"+f.Name()+"/", -1)
		re := regexp.MustCompile("(.+?);")
		ds := re.FindStringSubmatch(dataSource)
		dataSourceDev = append(dataSourceDev, ds[1])
		log.Println("  dev:", ds[1])

		data, _ = ioutil.ReadFile(path + "/" + f.Name() + "/src/main/resources/config/application-prod.yml")
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource2 := m["spring"].(map[interface{}]interface{})["datasource"].(map[interface{}]interface{})["url"].(string)
		re = regexp.MustCompile("(.+?)\\?")
		ds = re.FindStringSubmatch(dataSource2)
		dataSourceProd = append(dataSourceProd, ds[1])
		log.Println("  prod:", ds[1])

		data, _ = ioutil.ReadFile(path + "/docker-compose/docker-compose.yml")
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource3 := m["services"].(map[interface{}]interface{})[strings.ToLower(f.Name())+"-app"].(map[interface{}]interface{})["environment"].([]interface{})
		dataSource3Formatted := make([]string, len(dataSource3))
		for i, v := range dataSource3 {
			dataSource3Formatted[i] = fmt.Sprint(v)
		}
		re = regexp.MustCompile("SPRING_DATASOURCE_URL=(.+?)\\?")
		for _, dst := range dataSource3Formatted {
			ds = re.FindStringSubmatch(dst)
			if ds != nil {
				break
			}
		}
		dataSourceDocker = append(dataSourceDocker, ds[1])
		log.Println("  docker:", ds[1])
	}

	return dataSourceDev, dataSourceProd, dataSourceDocker
}
