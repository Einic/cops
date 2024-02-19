/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: tools
 * @Version: 1.0.0
 * @Date: 2024/2/18 17:42
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package utils

import (
	"context"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"github.com/Einic/cops/lib"
	AlterResource "github.com/Einic/cops/resources"
	"github.com/Einic/cops/zaplog"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"unicode"
)

func PrintVersionAndMD5() {
	// Get MD5 hash of the executable
	md5Hash, err := calculateMD5(os.Args[0])
	if err != nil {
		fmt.Println("Error calculating MD5 hash:", err)
		return
	}
	printDefaultContent()
	fmt.Printf("cdmasst version %s\n", lib.Version)
	fmt.Printf("MD5 hash: %x\n", md5Hash)
}

func calculateMD5(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func printDefaultContent() {
	// Default content when no valid flags are provided
	fmt.Println(`
***********************************************
* @Author: Einic <einicyeo AT gmail.com>      *
* @Description:                               *
* @BLOG:  https://www.infvie.com              *
* @Project home page:                         *
*     @https://github.com/Einic/EnvoyinStack  *
***********************************************`)
}

// ParseCSV parses a CSV file and returns its content as a 2D slice.
func ParseCSV(csvPath string) ([][]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Read and discard the header line
	if err != nil {
		return nil, err
	}

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// ValidateFields checks if any field in a CSV line is empty.
func ValidateFields(line []string) bool {
	for _, field := range line {
		if field == "" {
			return false
		}
	}
	return true
}

// IsMilliCPU checks if a CPU limit/request value is in milli-units.
func IsMilliCPU(cpu string) bool {
	if strings.HasSuffix(cpu, "m") && len(cpu) > 1 && unicode.IsDigit(rune(cpu[len(cpu)-2])) {
		return true
	}
	return false
}

// IsMegaMemory checks if a memory limit/request value is in Mebibytes.
func IsMegaMemory(memory string) bool {
	if strings.HasSuffix(memory, "Mi") && len(memory) > 2 && unicode.IsDigit(rune(memory[len(memory)-3])) {
		return true
	}
	return false
}

// UpdateWorkload updates the specified workload based on its type.
func UpdateWorkload(clientset *kubernetes.Clientset, worktype, namespace, workload, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory string, replicas int, logger zaplog.Logger) (lib.ResourceInfo, error) {
	var update lib.ResourceInfo

	switch worktype {
	case "deployment":
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), workload, metav1.GetOptions{})
		if err != nil {
			return update, fmt.Errorf("error getting deployment %s in namespace %s: %v", workload, namespace, err)
		}
		update = AlterResource.UpdateDeployment(clientset, deployment, replicas, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace, logger)

	case "statefulset":
		statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), workload, metav1.GetOptions{})
		if err != nil {
			return update, fmt.Errorf("error getting statefulset %s in namespace %s: %v", workload, namespace, err)
		}
		update = AlterResource.UpdateStatefulSet(clientset, statefulSet, replicas, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace, logger)

	default:
		return update, fmt.Errorf("unsupported worktype: %s", worktype)
	}

	return update, nil
}
