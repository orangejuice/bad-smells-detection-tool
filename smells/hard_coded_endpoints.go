package smells

import (
	"../util"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func ScanEndpoints(services []os.FileInfo, path string) (bool, []string) {
	log.Println("task -> scan hard-coded endpoints")

	var hardCoded bool
	var affectedService []string

	for _, f := range services {
		var dev, prod, docker bool

		log.Println("scanning", f.Name())

		// eureka enabled
		data, _ := ioutil.ReadFile(path + "/" + f.Name() + "/src/main/resources/config/application-dev.yml")
		m := make(map[interface{}]interface{})
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource := m["eureka"].(map[interface{}]interface{})["client"].(map[interface{}]interface{})["service-url"].(map[interface{}]interface{})["defaultZone"].(string)
		re := regexp.MustCompile("@((.+?):.+?)/")
		ds := re.FindStringSubmatch(dataSource)
		if !dev && util.IsIpAddress(ds[2]) {
			dev = true
		}
		log.Println("  dev:", "  -", "eureka address:", ds[1])

		data, _ = ioutil.ReadFile(path + "/" + f.Name() + "/src/main/resources/config/application-prod.yml")
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource2 := m["eureka"].(map[interface{}]interface{})["client"].(map[interface{}]interface{})["service-url"].(map[interface{}]interface{})["defaultZone"].(string)
		re = regexp.MustCompile("@((.+?):.+?)/")
		ds = re.FindStringSubmatch(dataSource2)
		if !prod && util.IsIpAddress(ds[2]) {
			prod = true
		}
		log.Println("  prod:", "-", "eureka address:", ds[1])

		data, _ = ioutil.ReadFile(path + "/docker-compose/docker-compose.yml")
		_ = yaml.Unmarshal([]byte(data), &m)
		dataSource3 := m["services"].(map[interface{}]interface{})[strings.ToLower(f.Name())+"-app"].(map[interface{}]interface{})["environment"].([]interface{})
		dataSource3Formatted := make([]string, len(dataSource3))
		for i, v := range dataSource3 {
			dataSource3Formatted[i] = fmt.Sprint(v)
		}
		re = regexp.MustCompile("EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE=.+?@((.+?):.+?)/")
		for _, dst := range dataSource3Formatted {
			ds = re.FindStringSubmatch(dst)
			if ds != nil {
				break
			}
		}
		if !docker && util.IsIpAddress(ds[2]) {
			docker = true
		}
		log.Println("  docker:", "  -", "eureka address:", ds[1])

		if dev || prod || docker {
			hardCoded = true
			affectedService = append(affectedService, f.Name())
		}
	}
	return hardCoded, affectedService
}
