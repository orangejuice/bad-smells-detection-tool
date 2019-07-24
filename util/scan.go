package util

import (
	"io/ioutil"
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
					break
				}
			}
		}
	}
	return services
}
