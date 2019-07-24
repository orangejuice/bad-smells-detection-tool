package main

import (
	"./smells"
	"./util"
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	f := util.InitLog()
	defer f.Close()

	var path string

	fmt.Println(os.Args)
	if len(os.Args) == 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter project path: ")
		path, _ = reader.ReadString('\n')

	} else {
		path = os.Args[1]
	}

	fmt.Println("scanning for micro-services in", path)
	time.Sleep(1 * time.Second)

	//filter, only leaves micro-services
	services := util.ScanMicroservices(path)

	fmt.Println("discovered micro-services:")
	for _, f := range services {
		fmt.Println("-", f.Name())
	}

	fmt.Println("\ndetecting API gateway")
	time.Sleep(1 * time.Second)

	fmt.Println("\nAPI Gateway assessment report")
	fmt.Println("============================")
	NAGDetected, NAGServices := smells.ScanApiGateway(services, path)
	fmt.Println("  gateway detected:", NAGDetected)
	if NAGDetected {
		fmt.Println("  in services:", NAGServices)
	}
	fmt.Println("============================")

	fmt.Println("\ndetecting data management configurations")
	time.Sleep(1 * time.Second)

	dataSourceDev, dataSourceProd, dataSourceDocker := smells.ScanSDM(services, path)

	fmt.Println("\nShared Data Management report")
	fmt.Println("============================")
	fmt.Println("  2 spring profiles and 1 docker-compose config")
	fmt.Println("  dev", "-", smells.IsSharedDataManagement(dataSourceDev))
	fmt.Println("  prod", "-", smells.IsSharedDataManagement(dataSourceProd))
	fmt.Println("  docker", "-", smells.IsSharedDataManagement(dataSourceDocker))
	fmt.Println("============================")

	fmt.Println("\ndetecting hard-coded endpoints")
	time.Sleep(1 * time.Second)

	HCEDetected, HCEServices := smells.ScanEndpoints(services, path)

	fmt.Println("\nHard-Coded Endpoints report")
	fmt.Println("============================")
	fmt.Println("  existence of risk:", HCEDetected)
	fmt.Println("  in services:", HCEServices)
	fmt.Println("============================")

	fmt.Println("\ndetecting circuit breakers")
	time.Sleep(1 * time.Second)

	fmt.Println("\nCircuit Breaker assessment report")
	fmt.Println("============================")
	NCBRisk, NCBRiskServices := smells.ScanCircuitBreaker(services, path)
	fmt.Println("  existence of risk:", NCBRisk)
	if NCBRisk {
		fmt.Println("  in services:", NCBRiskServices)
	}
	fmt.Println("============================")

	fmt.Println("\ndetails are available in", f.Name(), "file")
}
