package util

import (
	"io/ioutil"
)

func ScanConfigurationProfiles(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var configs []string

	for _, f := range files {
		if f.IsDir() {
			files, _ := ioutil.ReadDir(path + "/src/main/resources/config")

			for _, ft := range files {
				if !ft.IsDir() {
					configs = append(configs, ft.Name())
					break
				}
			}
		}
	}
	return configs
}

func ReadYamlFile(t interface{}, path string) {

}
