package main

import (
	"fmt"
	"io/ioutil"
	"reflect" //for reading out type of variable
)


func TestAddLowSpeedStates_1(t *testing.T) {
	//Read File
	dat, err := ioutil.ReadFile("C:/Users/LarsWitte/umh_larswitte_repo/sensorconnectRep/united-manufacturing-hub/golang/cmd/sensorconnect/ifm-0002BA-20170227-IODD1.1.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println("Contents of file:", string(dat))

	// Unmarshal file
	var ioDevice IoDevice
	ioDevice, err = UnmarshalIoddFile(dat)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	//should give out 698 and int
	if !reflect.DeepEqual(ioDevice.ProfileBody.DeviceIdentity.DeviceId), 698) {
		t.Error()
	}

	if !reflect.DeepEqual(reflect.TypeOf(ioDevice.ProfileBody.DeviceIdentity.DeviceId), "int") {
		t.Error()
	}
	


	//should put out 4 and int
	fmt.Println(ioDevice.ProfileBody.DeviceFunction.ProcessDataCollection.ProcessData.ProcessDataIn.Datatype.ReccordItem[1].SimpleDatatype.BitLength)
	fmt.Println(reflect.TypeOf(ioDevice.ProfileBody.DeviceFunction.ProcessDataCollection.ProcessData.ProcessDataIn.Datatype.ReccordItem[1].SimpleDatatype.BitLength))
	fmt.Println(reflect.TypeOf(dat))

}