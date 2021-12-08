package main

import (
	"fmt"
	"io/ioutil"
	"reflect" //for reading out type of variable
	"testing"
)

func TestUnmarshalIoddFile(t *testing.T) {
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

	//DeviceId: should give out 698
	if !reflect.DeepEqual(ioDevice.ProfileBody.DeviceIdentity.DeviceId, 698) {
		t.Error()
	}
	//DeviceId: type should be int
	if !reflect.DeepEqual(reflect.TypeOf(ioDevice.ProfileBody.DeviceIdentity.DeviceId).Kind(), reflect.Int) {
		t.Error()
	}

	//BitLength (of SimpleDatatype): should be 4 here
	if !reflect.DeepEqual(ioDevice.ProfileBody.DeviceFunction.ProcessDataCollection.ProcessData.ProcessDataIn.Datatype.ReccordItem[1].SimpleDatatype.BitLength, 4) {
		t.Error()
	}
	//BitLength (of SimpleDatatype): type should be int
	if !reflect.DeepEqual(reflect.TypeOf(ioDevice.ProfileBody.DeviceFunction.ProcessDataCollection.ProcessData.ProcessDataIn.Datatype.ReccordItem[1].SimpleDatatype.BitLength).Kind(), reflect.Int) {
		t.Error()
	}
}
