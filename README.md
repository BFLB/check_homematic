# HomeMatic

Golang Monitoring plugin for homematic devices 

First basic implementation

Feedback, testing and contribution welcome!

Supported device-types:
- HM-WDS40-TH-I-2
- HM-WDS10-TH-O

Usage:
- check_homematic -H=\<Host\> -devType=<DeviceType> -devName=<DeviceName> [-wTemp=<range>] [-cTemp=<range>] [-wHumi=<range>] [-cHumi=<range>] 

Arguments:
- H string HomeMatic Address (mandatory)
- devName string Device name (mandatory)
- devType string Device type (mandatory)
- p string default "80" Port (optional)
- wTemp string Warning Threshold Temperature (optional)
- cTemp string Critical Threshold Temperature (optional)
- Humi string Warning Threshold Humidity (optional)
- cHumi string Critical Threshold Humidity (optional)

