package smells

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanApiGateway(services []os.FileInfo, path string) (bool, []string) {
	var detected bool
	var serviceName []string
	n := 0
	log.Println("task -> scan api gateway")

	for _, service := range services {
		log.Println("scanning", service.Name())

		var files []string
		_ = filepath.Walk(path+"/"+service.Name()+"/src/main/java", func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".java" {
				files = append(files, path)
			}
			return nil
		})
		log.Println("java code files ready to be scanned:", len(files))

		for _, file := range files {
			data, _ := ioutil.ReadFile(file)
			n++
			if strings.Contains(string(data), "@EnableZuulProxy") {
				fmt.Println(" ", n, "files scanned")
				log.Println("found Gateway")
				detected = true
				serviceName = append(serviceName, service.Name())
				break
			}
		}
	}

	return detected, serviceName
}
