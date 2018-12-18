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

	hm           "github.com/BFLB/HomeMatic"
	check        "github.com/BFLB/monitoringplugin"
	r            "github.com/BFLB/monitoringplugin/Range"
	activeWriter "github.com/BFLB/monitoringplugin/writers/activeWriter"
	hmwds40thi   "github.com/BFLB/check_homematic/devices/hmwdsfortythi"

)

// Comman-line Arguments
var (
	host          = flag.String("H", "", "HomeMatic Address (mandatory)")
	devName       = flag.String("devName", "", "Device name (mandatory)")
	devType       = flag.String("devType", "", "Device type (mandatory)")
	port          = flag.String("p", "80", "Port (optional)")
	wTemp         = flag.String("wTemp", "", "Warning Threshold Temperature (optional)")
	cTemp         = flag.String("cTemp", "", "Critical Threshold Temperature (optional)")
	wHumi         = flag.String("wHumi", "", "Warning Threshold Humidity (optional)")
	cHumi         = flag.String("cHumi", "", "Critical Threshold Humidity (optional)")
)

func main() {

	
	// Create new check
	check   := check.New()
	message := ""

	// Create writer
	writer := activeWriter.New()
	
	
	// Fixme
	// Print usage info (Override of flag.Usage)
	flag.Usage = func() {
		check.Message("Usage:")
		check.Message("  -H string HomeMatic Address (mandatory)")
		check.Message("  -devName string Device name (mandatory)")
		check.Message("  -devType string Device type (mandatory)")
		check.Message("  -p string default \"80\" Port (optional)")
		check.Message("  -wTemp string Warning Threshold Temperature (optional)")
		check.Message("  -cTemp string Critical Threshold Temperature (optional)")
		check.Message("  -Humi string Warning Threshold Humidity (optional)")
		check.Message("  -cHumi string Critical Threshold Humidity (optional)")
		check.Status.Unknown()
		writer.Write(check)		
	}
	
	// Parse command-line args
	flag.Parse()

	// Check mandatory args
	if *host == "" {
		flag.Usage()	
	}
	if *devName == "" {
		flag.Usage()	
	}
	if *devType == "" {
		flag.Usage()	
	}

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

	// Init HomeMatic connection
	api, err := hm.Init(*host, *port)
	if err != nil {
		log.Fatal("Init returned error: ", err)
	}

	// Get device-id
	deviceID   := ""

	// Get devicelist
	dl, err := api.DeviceList()
	if err != nil {
		check.Status.Critical(true)
		check.Message(err.Error())
		writer.Write(check)
	}

	// Get device by name and type
	for i := 0; i < len(dl.Device); i++ {
		if dl.Device[i].Name == *devName && dl.Device[i].DeviceType == *devType {
			deviceID = dl.Device[i].IseID
			break
		}
	}

	// Device not found error
	if deviceID == "" {
		message = fmt.Sprintf("No device with name=%s and type=%s found", *devName, *devType)
		check.Status.Critical(true)
		check.Message(message)
		writer.Write(check)
	}

	// Get state
	state, err := api.State(deviceID)
	if err != nil {
		check.Status.Critical(true)
		check.Message(err.Error())
		writer.Write(check)
	}

	// Run specific checks
	switch *devType {
	case "HM-WDS40-TH-I-2", "HM-WDS10-TH-O":
		hmwds40thi.Check(state, wTemp, cTemp, wHumi, cHumi, writer)

	default:
		message = fmt.Sprintf("Check for device type =%s missing", *devType)
		check.Status.Critical(true)
		check.Message(message)
		writer.Write(check)
	}

	// Unknown error
	message = fmt.Sprintf("Unknown error")
	check.Status.Critical(true)
	check.Message(message)
	writer.Write(check)

}
