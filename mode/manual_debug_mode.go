/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: manual_debug_mode
 * @Version: 1.0.0
 * @Date: 2024/2/19 12:49
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package mode

import (
	"fmt"
	"github.com/Einic/cops/zaplog"
	"github.com/guptarohit/asciigraph"
	"math/rand"
	"time"
)

func ManualDebugMode(logger zaplog.Logger) {
	logger.Debug("ManualDebugMode")
}

func PrintDebugGraph() {
	rand.Seed(time.Now().UnixNano())
	var data []float64
	for i := 0; i < 160; i++ {
		price := rand.Float64() * 100
		data = append(data, price)
	}
	graph := asciigraph.Plot(data, asciigraph.Height(6),
		asciigraph.Caption("@EinicYeo"),
		asciigraph.CaptionColor(asciigraph.Green),
		asciigraph.AxisColor(asciigraph.MediumSlateBlue),
		asciigraph.LabelColor(asciigraph.LightGray),
		asciigraph.SeriesColors(asciigraph.DodgerBlue),
	)
	fmt.Println(graph)
}
