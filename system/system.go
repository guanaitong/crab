package system

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const unknown = "unknown"

var (
	inK8s                bool
	appName, appInstance string
	workEnv, workIdc     string
	hostName, hostIp     string
	instanceId           string
	startTime            time.Time
)

func SetupAppName(appN string) {
	appName = appN
	log.Printf("setupAppName %s", appName)
}

func init() {
	inK8s = len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0
	if inK8s {
		log.Printf("app is running inside k8s")
	} else {
		log.Printf("app is running outside k8s")
	}

	appName = getEnv("APP_NAME", unknown)

	//在k8s里，使用HOSTNAME，VM里使用APP_INSTANCE_NAME
	if inK8s {
		appInstance = getEnv("HOSTNAME", unknown)
	} else {
		appInstance = getEnv("APP_INSTANCE_NAME", unknown)
	}

	log.Printf("appName %s, appInstance %s", appName, appInstance)

	workEnv, workIdc = getEnv("WORK_ENV", "dev"), getEnv("WORK_IDC", "ofc")

	log.Printf("workEnv %s, workIdc %s", workEnv, workIdc)

	if inK8s {
		hostName = os.Getenv("HOSTNAME")
	} else {
		hostName, _ = os.Hostname()
	}

	hostIp = getLocalIP()

	log.Printf("hostName %s, hostIp %s", hostName, hostIp)

	instanceId = appName + "-->" + appInstance + "-->" + hostIp + "-->" + fmt.Sprint(rand.Int63n(time.Now().UnixNano()))

	log.Printf("instanceId %s", instanceId)

	startTime = time.Now()

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func InK8s() bool {
	return inK8s
}

func GetAppName() string {
	return appName
}

func GetAppInstance() string {
	return appInstance
}

func GetWorkEnv() string {
	return workEnv
}
func GetWorkIdc() string {
	return workIdc
}

func GetHostName() string {
	return hostName
}

func GetHostIp() string {
	return hostIp
}

func GetInstanceId() string {
	return instanceId
}

func GetStartTime() time.Time {
	return startTime
}

func GetServiceDomainSuffix() string {
	if workEnv == "prepare" && workIdc == "sh" {
		return ".services.product.sh"
	}
	return ".services." + workEnv + "." + workIdc
}

func GetServerDomainSuffix() string {
	if workEnv == "prepare" && workIdc == "sh" {
		return ".server.product.sh"
	}
	return ".server." + workEnv + "." + workIdc
}
