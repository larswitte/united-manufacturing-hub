---
title: "Sensorconnect"
linkTitle: "Sensorconnect"
description: >
  This microservice of the united manufacturing hub automatically detects ifm gateways in the specified network and reads their sensor values in the  in the highest possible data frequency (or chosen data frequency). It preprocesses the data and sends it via [MQTT](/docs/concepts/mqtt/) or Kafka.
aliases:
  - /docs/Developers/factorycube-edge/sensorconnect
  - /docs/developers/factorycube-edge/sensorconnect
---
## Sensorconnect overview
Sensorconnect provides plug-and-play access to [IO-Link](https://io-link.com/en/) sensors connected to [ifm gateways](https://www.ifm.com/us/en/category/245_010_010). Digital input mode is also supported. A typical setup contains multiple sensors connected to ifm gateways on the production shopfloor. Those gateways are connected via LAN to your server infrastructure. The sensorconnect microservice constantly monitors a given IP range for gateways. Once a gateway is found, it automatically starts requesting, receiving and processing sensordata in short intervals. The received data is preprocessed based on a database including thousands of sensor definitions. And because we strongly believe in open industry standards, Sensorconnect brings the data via MQTT or Kafka to your preferred software solutions, for example, features of the United Manufacturing Hub or the cloud.

## Which problems is Sensorconnect solving
Let's take a step back and think about, why we need a special microservice:
1. The gateways are providing a rest api to request sensordata. Meaning as long as we dont have a process to request the data, nothing will happen. 
2. Constantly requesting and processing data with high robustness and reliability can get difficult in setups with a large number of sensors.
3. Even if we use for example node-red flows ([flow from node-red](https://flows.nodered.org/node/sense-ifm-iolink), [flow from ifm](https://www.ifm.com/na/en/shared/technologies/io-link/system-integration/iiot-integration)) to automatically request data from the rest api of the ifm gateways, the information is most of the times cryptic without proper interpretation. 
4. Device manufacturers will provide one IODD file (IO Device Description), for every sensor and actuator they produce. Those contain information to correctly interpret data from the devices. They are in XML-format and available at the [IODDfinder website](https://io-link.com/en/IODDfinder/IODDfinder.php). But automatic IODD interpretation is relatively complex and manually using IODD files is not economically feasible.

## Installation
### Production setup
Sensorconnect comes directly with the united manufacturing hub - no additional installation steps required. We allow the user to customize sensorconnect by changing environment variables. All possible environment variables you can use to customize sensorconnect and how you change them is described in the [environment-variables documentation](/docs/Developers/united-manufacturing-hub/environment-variables/). Sensorconnect is by default enabled in the United Manufacturing Hub. To set your preferred serialnumber and choose the ip range to look for gateways, you can configure either the values.yaml directly or use our managment SaaS tool. 

### Development setup
Here is a quick tutorial on how to start up a basic configuration / a basic docker-compose stack, so that you can develop.

1. execute `docker-compose -f ./deployment/sensorconnect/docker-compose.yaml up -d --build`


## Underlying Functionality
Sensorconnect downloads relevant IODD files automatically after installation from the [IODDfinder website](https://io-link.com/en/IODDfinder/IODDfinder.php). If an unknown sensor is connected later, sensorconnect will automatically download the file. We will also provide a folder to manually deposit IODD-files, if the automatic download doesn't work (e.g. no internet connection).

### Rest API POST requests from sensorconnect to the gateways
Sensorconnect scans the ip range for new ifm gateways (used to connect the IO-Link devices to). To do that, sensorconnect iterates through all the possible IP addresses in the specified IP address range ("http://"+url, `payload`, timeout=0.1). It stores the IP adresses, with the product codes (*the types of the gateways*) and the individual serialnumbers.

**Scanning with following `payload`: (information sent during a POST request to the ifm gateways)**
```JSON
{
  "code":"request",
  "cid":-1, // The cid (Client ID) can be chosen.
  "adr":"/getdatamulti",
  "data":{
    "datatosend":[
      "/deviceinfo/serialnumber/","/deviceinfo/productcode/"
      ]
    }
}
```
**Example answer from gateway:**
```JSON
{
    "cid": 24,
    "data": {
        "/deviceinfo/serialnumber/": {
            "code": 200,
            "data": "000201610192"
        },
        "/deviceinfo/productcode/": {
            "code": 200,
            "data": "AL1350"
        }
    },
    "code": 200
}
```
All port modes of the connected gateways are requested. Depending on the productcode, we can determine the total number of ports on the gateway and iterate through them.

**Requesting port modes with following payload: (information sent during a POST request to the ifm gateways)**
```JSON
{
  "code":"request",
  "cid":-1,
  "adr":"/getdatamulti",
  "data":{
    "datatosend":[
      "/iolinkmaster/port[1]/mode",
      "/iolinkmaster/port<i>/mode" //looping through all ports on gateway
      ]
    }
}
```
**Example answer from gateway:**
```JSON
{
    "cid": -1,
    "data": {
        "/iolinkmaster/port[1]/mode": {
            "code": 200,
            "data": 3
        },
        "/iolinkmaster/port[2]/mode": {
            "code": 200,
            "data": 3
        }
    }
}
```
If the mode == 1: port_mode = "DI" (Digital Input)
If the mode == 2: port_mode = "DO" (Digital output)
If the mode == 3: port_mode = "IO_Link"

All values of accessible ports are requested as fast as possible (ifm gateways are by far the bottleneck in comparison to the networking).
**Requesting IO_Link port values with following payload: (information sent during a POST request to the ifm gateways)**
```JSON
{
  "code":"request",
  "cid":-1,
  "adr":"/getdatamulti",
  "data":{
    "datatosend":[
      "/iolinkmaster/port[1]/iolinkdevice/deviceid",
      "/iolinkmaster/port[1]/iolinkdevice/pdin",
      "/iolinkmaster/port[1]/iolinkdevice/vendorid",
      "/iolinkmaster/port[1]/pin2in",
      "/iolinkmaster/port[<i>]/iolinkdevice/deviceid",//looping through all connected ports on gateway
      "/iolinkmaster/port[<i>]/iolinkdevice/pdin",
      "/iolinkmaster/port[<i>]/iolinkdevice/vendorid",
      "/iolinkmaster/port[<i>]/pin2in"
      ]
  }
}
```
**Example answer from gateway:**
```JSON
{
    "cid": -1,
    "data": {
        "/iolinkmaster/port[1]/iolinkdevice/deviceid": {
            "code": 200,
            "data": 278531
        },
        "/iolinkmaster/port[1]/iolinkdevice/pdin": {
            "code": 200,
            "data": "0101" // This string contains the actual data from the sensor. In this example it is the data of a buttonbar. The value changes when one button is pressed. Interpreting this value automatically relies heavily on the IODD file of the specific sensor (information about which partition of the string holds what value (type, value name etc.)).
        },
        "/iolinkmaster/port[1]/iolinkdevice/vendorid": {
            "code": 200,
            "data": 42
        },
        "/iolinkmaster/port[1]/pin2in": {
            "code": 200,
            "data": 0
        }
    },
    "code": 200
}
```

### IODD file management
Based on the vendorid and deviceid (extracted out of the received data from the ifm-gateway), sensorconnect looks in it's persistant storage if an iodd file is available. If there is no iodd file stored for the specific sensor, it tries to download a file online and saves it. If this also is not possible, sensorconnect doesn't preprocess the pdin data entry and forwards it as is.

### Data interpretation (key: pdin) with IODD files
The iodd files are in xml format and often contain multiple thousand lines. You can find [extensive documentation on the io link website](https://io-link.com/de/Download/Download.php). Especially relevant areas of the iodd files are (for our use-case):
1. IODevice/DocumentInfo/version: version of IODD files
2. IODevice/ProfileBody/DeviceIdentity: deviceid, vendorid etc.
3. IODevice/DeviceFunction/ProcessDataCollection/ProcessData/ProcessDataIn: contains information about the data structure received from the sensor. Datatypes or datatype references are given.
4. IODevice/DeviceFunction/DatatypeCollection: if used, contains datatypes referenced by the ProcessDataCollection(Point 3.)
5. IODevice/ExternalTextCollection: contains translations for textId's.

### Digital input data management
The gateways are also able to receive digital input data (not IO-Link). This data is also requested and forwarded.

### Data delivery to MQTT or Kafka
For delivery of the data, sensorconnect converts it to a JSON containing the preprocessed data, timestamps, serialnumber etc. and sends it via MQTT to the MQTT broker or via Kafka to the Kafka broker. You can change the [environment variables](/docs/Developers/united-manufacturing-hub/environment-variables/) to choose your preferred protocol. The format of those JSON messages coming from sensorconnect is described in detail and with examples on the [UMH Datamodel website](/docs/concepts/mqtt/).


## Tested gateways
AL1350
AL1342

## Future plans
- Support for additional gateway manufacturers.
- Support for value ranges and single values specified in the IODD file (giving additional information about currently received sensor value).