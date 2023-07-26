package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
	"time"
)

var stat unix.Statfs_t

func main() {
	monitoredPath := os.Getenv("MONITORED_PATH")
	if monitoredPath == "" {
		monitoredPath = "/"
	}
	go func() {
		http.Handle("/", promhttp.Handler())
		err := http.ListenAndServe(":1971", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// get
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			getDiskUsage(monitoredPath)
		}
	}
}

func getDiskUsage(path string) {
	err := unix.Statfs(path, &stat)
	if err != nil {
		fmt.Println(err)
		return
	}
	bAvail := stat.Bavail * uint64(stat.Bsize)
	bSize := stat.Blocks * uint64(stat.Bsize)
	bUsed := bSize - bAvail
	bPercent := float64(bUsed) / float64(bSize) * 100

	metricUsed.WithLabelValues(path).Set(float64(bUsed))
	metricPercent.WithLabelValues(path).Set(bPercent)
	metricAvail.WithLabelValues(path).Set(float64(bAvail))
}

var metricUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "total_disk_used",
	Help: "Disk used",
}, []string{"path"})

var metricPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "total_disk_percent",
	Help: "Disk usage percent",
}, []string{"path"})

var metricAvail = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "total_disk_avail",
	Help: "Disk available",
}, []string{"path"})
