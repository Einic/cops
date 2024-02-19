/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: normal_mode
 * @Version: 1.0.0
 * @Date: 2024/2/19 12:49
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package mode

import (
	"flag"
	"fmt"
	"github.com/Einic/cops/lib"
	"github.com/Einic/cops/table"
	"github.com/Einic/cops/utils"
	"github.com/Einic/cops/zaplog"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strconv"
)

func NormalMode(logger zaplog.Logger) {
	// Define flags
	versionFlag := flag.Bool("v", false, "Print version number and MD5 hash")
	versionLongFlag := flag.Bool("version", false, "Print version number and MD5 hash")
	helpFlag := flag.Bool("h", false, "Show help message")
	helpLongFlag := flag.Bool("help", false, "Show help message")
	alterFlag := flag.String("a", "", "Please alter resource")
	alterLongFlag := flag.String("alter", "", "Please alter resource")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Printf("  -v, --version   Print version number and MD5 hash.\n")
		fmt.Printf("  -h, --help      Please read README.md to configure.\n")
		fmt.Printf("  -a, --alter     Please alter resource [-a /root/.kube/config ./example.csv].\n")
	}

	// Parse flags
	flag.Parse()

	// Check flags and execute corresponding actions
	if *versionFlag || *versionLongFlag {
		utils.PrintVersionAndMD5()
	} else if *helpFlag || *helpLongFlag {
		flag.Usage()
	} else if *alterFlag != "" || *alterLongFlag != "" {
		args := os.Args[2:]
		executeCommand(logger, args...)
	} else {
		// If an unknown flag is provided, print custom message
		fmt.Printf("Unknown flag provided: %s\n", os.Args[1:])
		flag.Usage()
	}
}

func executeCommand(logger zaplog.Logger, args ...string) {
	if len(args) != 2 {
		logger.Warn("Invalid number of arguments for alter resource. Expected 2, got ", zap.Int("LenArgs", len(args)))
		flag.Usage()
		return
	}

	lib.Kubeconfig, lib.CSVPath = args[0], args[1]

	config, err := clientcmd.BuildConfigFromFlags("", lib.Kubeconfig)
	if err != nil {
		logger.Error("Error building kubeconfig", zap.Error(err))
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("Error creating clientset", zap.Error(err))
		os.Exit(1)
	}

	lines, err := utils.ParseCSV(lib.CSVPath)
	if err != nil {
		logger.Error("Error parsing CSV file", zap.Error(err))
		os.Exit(1)
	}

	var updates []lib.ResourceInfo

	// Launch goroutines to handle each line of the CSV file
	for _, line := range lines {
		if len(line) != 9 {
			logger.Error("Invalid CSV format. Expected 9 fields per line.")
			continue
		}

		if !utils.ValidateFields(line) {
			logger.Error("Empty field found in CSV.")
			continue
		}

		// Parse fields from the CSV line
		workload := line[0]
		containersName := line[1]
		worktype := line[2]
		namespace := line[3]
		replicasCSV, err := strconv.Atoi(line[4])
		if err != nil {
			logger.Error("Error converting replicas to integer", zap.Error(err))
			continue
		}
		limitsCPU := line[5]
		limitsMemory := line[6]
		requestsCPU := line[7]
		requestsMemory := line[8]

		// Update the workload based on worktype
		if !utils.IsMilliCPU(limitsCPU) || !utils.IsMilliCPU(requestsCPU) {
			logger.Error("CPU limit/request should be in milli-units (suffix 'm').", zap.String("Workload", workload), zap.String("Namespace", namespace))
			continue
		}

		if !utils.IsMegaMemory(limitsMemory) || !utils.IsMegaMemory(requestsMemory) {
			logger.Error("Memory limit/request should be in Mebibytes (suffix 'Mi').", zap.String("Workload", workload), zap.String("Namespace", namespace))
			continue
		}

		update, err := utils.UpdateWorkload(clientset, worktype, namespace, workload, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, replicasCSV, logger)
		if err != nil {
			logger.Error("Error updating workload", zap.String("Workload", workload), zap.String("Namespace", namespace), zap.Error(err))
			continue
		}
		updates = append(updates, update)
	}
	table.PrintUpdateTable(updates)
}
