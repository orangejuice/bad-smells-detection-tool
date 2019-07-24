package smells

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanCircuitBreaker(services []os.FileInfo, path string) (bool, []string) {
	var risk bool
	var riskServices []string
	tn := 0
	log.Println("task -> scan circuit breaker")

	for _, f := range services {
		n := 0
		log.Println("scanning", f.Name())

		var files []string
		_ = filepath.Walk(path+"/"+f.Name()+"/src/main/java", func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".java" {
				files = append(files, path)
			}
			return nil
		})
		log.Println("  java code files ready to be scanned:", len(files))

		foundFeignCall := false
		for _, file := range files {
			data, _ := ioutil.ReadFile(file)
			n++
			if strings.Contains(string(data), "@AuthorizedFeignClient") || strings.Contains(string(data), "@FeignClient") {
				foundFeignCall = true
			}
		}
		log.Println(" ", n, "files scanned")
		tn += n
		// look for @feignclient | @AuthorizedFeignClient java file, should have feign:hystrix:enabled: true

		if foundFeignCall {
			data, _ := ioutil.ReadFile(path + "/" + f.Name() + "/src/main/resources/config/application.yml")
			m := make(map[interface{}]interface{})
			_ = yaml.Unmarshal([]byte(data), &m)
			enabled := m["feign"].(map[interface{}]interface{})["hystrix"].(map[interface{}]interface{})["enabled"].(bool)

			risk = !enabled
			log.Println("  found feign call, Hystrix circuit breaker enabled:", enabled)
			if risk {
				riskServices = append(riskServices, f.Name())
			}
		} else {
			log.Println("  no Feign Call detected")
		}
	}
	fmt.Println(" ", tn, "files scanned")
	return risk, riskServices
}
