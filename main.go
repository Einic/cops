/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2024/2/7 11:49
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package main

import (
	"fmt"
	"github.com/Einic/cops/lib"
	"github.com/Einic/cops/mode"
	"github.com/Einic/cops/zaplog"
)

const debugMode = false

func main() {
	logger := zaplog.InitLogger(1)
	defer logger.Close()

	if debugMode {
		fmt.Println("Debug mode is enabled, current version:", lib.Version)
		mode.PrintDebugGraph()
		mode.ManualDebugMode(logger)
	} else {
		mode.NormalMode(logger)
	}

}
