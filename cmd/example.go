package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/mem"

	"github.com/ariefrahmansyah/influxdb"
)

var hostname string
var influxClient influxdb.Client
var influxBatch influxdb.Batch

func main() {
	hostname, _ = os.Hostname()

	initInfluxClient()
	initInfluxBatch()

	go func() {
		for {
			if err := monitorMemory(); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	waitEnd()
}

func initInfluxClient() {
	// Make new http client connection
	infClient, err := influxdb.NewClient(influxdb.ClientConfig{
		Name:    "learn",
		Address: "http://localhost:8086",
		Type:    "http",
	})
	if err != nil {
		log.Fatalln(err)
	}
	influxClient = infClient
}

func initInfluxBatch() {
	// Make new batch points
	infBatch, err := influxdb.NewBatch(influxdb.BatchConfig{
		Database:  "mydb",
		Precision: "s",
	})
	if err != nil {
		log.Fatalln(err)
	}
	influxBatch = infBatch
}

func monitorMemory() error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	// Prepare the data
	tags := map[string]string{"hostname": hostname}
	fields := map[string]interface{}{
		"total_memory": fmt.Sprint(memory.Total),
		"free_memory":  fmt.Sprint(memory.Free),
		"used_memory":  fmt.Sprint(memory.UsedPercent),
	}

	// Make new point
	pt, err := client.NewPoint("memory", tags, fields, time.Now())
	if err != nil {
		return err
	}

	// Add point to batch points
	influxBatch.BP.AddPoint(pt)

	// Write point to batch points
	err = influxClient.Write(influxBatch)
	if err != nil {
		return err
	}

	log.Println(fields)

	return nil
}

func waitEnd() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	<-done
}
