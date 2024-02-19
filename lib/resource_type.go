/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: cnasst_type
 * @Version: 1.0.0
 * @Date: 2024/2/7 14:06
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package lib

var (
	Version    = "1.0.0"
	Kubeconfig string
	CSVPath    string
)

type ResourceInfo struct {
	DataTime              string
	Workload              string
	ContainerName         string
	WorkType              string
	Namespace             string
	CurrentReplicas       int
	AlterReplicas         int
	CurrentLimitsCPU      string
	AlterLimitsCPU        string
	CurrentLimitsMemory   string
	AlterLimitsMemory     string
	CurrentRequestsCPU    string
	AlterRequestsCPU      string
	CurrentRequestsMemory string
	AlterRequestsMemory   string
	PodQos                string
	RunStatus             string
	AlterStatus           string
}
