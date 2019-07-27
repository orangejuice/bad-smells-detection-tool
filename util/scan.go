package util

import (
	"io/ioutil"
	"log"
	"os"
)

func ScanMicroservices(path string) []os.FileInfo {
	files, _ := ioutil.ReadDir(path)
	var services []os.FileInfo

	for _, f := range files {
		if f.IsDir() {
			files, _ := ioutil.ReadDir(path + "/" + f.Name())

			for _, ft := range files {
				if ft.Name() == "src" {
					services = append(services, f)
					log.Println("Found microservices: ", f.Name())
					break
				}
			}
		}
	}
	log.Println("Found microservices number: ", len(services))
	return services
}
