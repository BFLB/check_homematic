// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Monitoring plugin for homematic CCU
//

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"os"

	hm       "github.com/BFLB/HomeMatic"
	check    "github.com/BFLB/monitoringplugin"
	status  "github.com/BFLB/monitoringplugin/status"
	perfdata "github.com/BFLB/monitoringplugin/performancedata"
	r        "github.com/BFLB/monitoringplugin/Range"
	writer   "github.com/BFLB/monitoringplugin/writers/activeWriter"
)

// Comman-line Arguments
var (
	host          = flag.String("H", "", "HomeMatic Address")
	port          = flag.String("p", "80", "Port")
	name          = flag.String("n", "", "Devicename")
)

func main() {

	// Parse command-line args
	flag.Parse()

	// Init HomeMatic connection
	api, err := hm.Init(*host, *port)
	if err != nil {
		log.Fatal("Init returned error: ", err)
	}

	// Check
	check   := check.New()
	message := ""
	rangeWarn := r.New()
	rangeCrit := r.New()

	rangeWarn.Parse("15:25")
	rangeCrit.Parse("10:30")


	// Get device-id
	deviceName := *name
	deviceID   := ""

	dl, err := api.DeviceList()
	if err != nil {
		check.Status.Critical(true)
		check.Message(err.Error())
		writer.Write(check)
	}

	for i := 0; i < len(dl.Device); i++ {
		if dl.Device[i].Name == deviceName {
			deviceID = dl.Device[i].IseID
			break
		}
	}

	state, err := api.State(deviceID)
	if err != nil {
		check.Status.Critical(true)
		check.Message(err.Error())
		writer.Write(check)
	}

	var unreach 	bool
	var lowBat  	bool
	var temperature float64
	var humidity	int64

	unreach = true
	unreach, _ = strconv.ParseBool(state.Device.Unreach) 

	lowBat = true
	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "LOWBAT" {
				lowBat, _ = strconv.ParseBool(state.Device.Channel[i].Datapoint[x].Value)
				break
			}
		}
	}

	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "TEMPERATURE" {
				temperature, _ = strconv.ParseFloat(state.Device.Channel[i].Datapoint[x].Value, 64)
				break
			}
		}
	}

	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "HUMIDITY" {
				humidity, _ = strconv.ParseInt(state.Device.Channel[i].Datapoint[x].Value, 10, 64)
				break
			}
		}
	}

	// Unreacheable
	if unreach == true {
		message += fmt.Sprintf("Device unreacheable ")
		check.Status.Critical(true)
	}

	// Lowbat
	if lowBat == true {
		message += fmt.Sprintf("Battery low ")
		check.Status.Warning(false)
	}

	// Temperature
	statTemp := status.New()
	statTemp.Threshold(temperature, &rangeWarn, &rangeCrit, true)
	check.Status.Merge(statTemp)

	message += fmt.Sprintf("Temperature: %.1f°C ", temperature)
	if statTemp.ReturnCode() != status.OK {
		message += fmt.Sprintf("(%s) ", statTemp.String())
	} 

	// Humidity
	statHumi := status.New()
	statHumi.Threshold(temperature, &rangeWarn, &rangeCrit, true)
	check.Status.Merge(statHumi)

	message += fmt.Sprintf("Humidity: %d%% ", humidity)
	if statHumi.ReturnCode() != status.OK {
		message += fmt.Sprintf("(%s) ", statTemp.String())
	} 

	// Add message
	check.Message(message)

	// Add Perfdata Temperature
	var datapoint *perfdata.PerformanceData

	datapoint, err = perfdata.New("Temperature", temperature, "", nil, nil, nil, nil )
	if err == nil {
		check.Perfdata(datapoint)
	}

	// Add Perfdata Humidity
	datapoint, err = perfdata.New("Humidity", float64(humidity), "%", nil, nil, nil, nil )
	if err == nil {
		check.Perfdata(datapoint)
	}

	fmt.Printf("%s", check.String())

	os.Exit(check.Status.ReturnCode())

}
