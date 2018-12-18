// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Monitoring plugin for homematic CCU
//

package hmwdsfortythi

import (
	"fmt"
	"strconv"

	hm       "github.com/BFLB/HomeMatic"
	check    "github.com/BFLB/monitoringplugin"
	s        "github.com/BFLB/monitoringplugin/status"
	perfdata "github.com/BFLB/monitoringplugin/performancedata"
	r        "github.com/BFLB/monitoringplugin/Range"
	writer   "github.com/BFLB/monitoringplugin/writers/activeWriter"
)


func Check(state *hm.State, wTemp *string, cTemp *string, wHumi *string, cHumi *string, w *writer.Writer) (status s.Status, err error) {

	// Create new check
	check   := check.New()
	message := ""
	
	// Set ranges
	var rangeWarnTemp *r.Range
	var rangeCritTemp *r.Range
	var rangeWarnHumi *r.Range
	var rangeCritHumi *r.Range

	if *wTemp != "" {
		rangeWarnTemp = r.New()
		rangeWarnTemp.Parse(*wTemp)
	}
	if *cTemp != "" {
		rangeCritTemp = r.New()
		rangeCritTemp.Parse(*cTemp)
	}
	if *wHumi != "" {
		rangeWarnHumi = r.New()
		rangeWarnHumi.Parse(*wHumi)
	}
	if *cHumi != "" {
		rangeCritHumi = r.New()
		rangeCritHumi.Parse(*cHumi)
	}

	// Variables
	var unreach 	bool
	var lowBat  	bool
	var temperature float64
	var humidity	int64

	// Check Unreach
	unreach = true
	unreach, _ = strconv.ParseBool(state.Device.Unreach) 

	if unreach == true {
		message += fmt.Sprintf("Device unreacheable ")
		check.Status.Critical(true)
	}

	// Check Lowbat
	lowBat = true
	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "LOWBAT" {
				lowBat, _ = strconv.ParseBool(state.Device.Channel[i].Datapoint[x].Value)
				break
			}
		}
	}

	if lowBat == true {
		message += fmt.Sprintf("Battery low ")
		check.Status.Warning(false)
	}

	// Check Temperature
	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "TEMPERATURE" {
				temperature, _ = strconv.ParseFloat(state.Device.Channel[i].Datapoint[x].Value, 64)
				break
			}
		}
	}

	statTemp := s.New()
	statTemp.Threshold(temperature, rangeWarnTemp, rangeCritTemp, true)
	check.Status.Merge(statTemp)

	message += fmt.Sprintf("Temperature: %.1f°C ", temperature)
	if statTemp.ReturnCode() != s.OK {
		message += fmt.Sprintf("(%s) ", statTemp.String())
	} 


	// Check Humidity
	for i := 0; i < len(state.Device.Channel); i++ {
		for x := 0; x < len(state.Device.Channel[i].Datapoint); x++ {
			if state.Device.Channel[i].Datapoint[x].Type == "HUMIDITY" {
				humidity, _ = strconv.ParseInt(state.Device.Channel[i].Datapoint[x].Value, 10, 64)
				break
			}
		}
	}

	statHumi := s.New()
	statHumi.Threshold(temperature, rangeWarnHumi, rangeCritHumi, true)
	check.Status.Merge(statHumi)

	message += fmt.Sprintf("Humidity: %d%% ", humidity)
	if statHumi.ReturnCode() != s.OK {
		message += fmt.Sprintf("(%s) ", statTemp.String())
	} 

	// Add message
	check.Message(message)

	// Add Perfdata Temperature
	var datapoint *perfdata.PerformanceData

	datapoint, err = perfdata.New("Temperature", temperature, "", rangeWarnTemp, rangeCritTemp, nil, nil )
	if err == nil {
		check.Perfdata(datapoint)
	}

	// Add Perfdata Humidity
	datapoint, err = perfdata.New("Humidity", float64(humidity), "%", rangeWarnHumi, rangeCritHumi, nil, nil )
	if err == nil {
		check.Perfdata(datapoint)
	}

	err = w.Write(check)

	return check.Status, err

}
