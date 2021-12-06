package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type IODevice struct {
	ProfileBody            ProfileBody `xml:"ProfileBody"`
	ExternalTextCollection `xml:"ExternalTextCollection"`
}

type ProfileBody struct {
	DeviceIdentity DeviceIdentity `xml:"DeviceIdentity"`
	DeviceFunction DeviceFunction `xml:"DeviceFunction"`
}

// DeviceIdentity : IODevice>ProfileBody>DeviceIdentity
type DeviceIdentity struct {
	VendorId   string `xml:"vendorId,attr"` // not used
	VendorName string `xml:"vendorName,attr"`
	DeviceId   string `xml:"deviceId,attr"` // id of type of device, given by device vendor
}

type ExternalTextCollection struct {
	PrimaryLanguage PrimaryLanguage `xml:"ExternalTextCollection"`
}

type PrimaryLanguage struct {
	Text []Text `xml:"Text"`
}

type Text struct {
	id    string `xml:"id,attr"`
	value string `xml:"value,attr"`
}

type DeviceFunction struct {
	ProcessDataCollection ProcessDataCollection `xml:"ProcessDataCollection"` //ToDo: array?
}

type ProcessDataCollection struct {
	ProcessData ProcessData `xml:"ProcessData"`
}

type ProcessData struct {
	ProcessDataIn ProcessDataIn `xml:"ProcessDataIn"`
}

type ProcessDataIn struct {
	Datatype Datatype
}

type Datatype struct {
	Bitlength   string     `xml:"bitLength,attr"`
	ReccordItem RecordItem `xml:"RecordItem"`
}

type RecordItem struct {
	BitOffset      int            `xml:"bitOffset,attr"`
	SimpleDatatype SimpleDatatype `xml:"SimpleDatatype"`
	Name           Name           `xml:"Name"`
}

type Name struct {
	TextId string `xml:"textId"`
}

type SimpleDatatype struct {
	Type        string `xml:"xsi:type,attr"` //ToDo how to unmarshal xsi:...
	BitLength   int    `xml:"bitLength,attr"`
	FixedLength int    `xml:"fixedLength,attr"`
}

func main() {
	payload := IODevice{}

	//Read File
	dat, err := ioutil.ReadFile("C:/Users/LarsWitte/umh_larswitte_repo/sensorconnectRep/united-manufacturing-hub/golang/cmd/sensorconnect/ifm-0002BA-20170227-IODD1.1.xml")
	check(err)
	fmt.Println("Contents of file:", string(dat))

	// Unmarshal file with Unmarshal
	err = xml.Unmarshal(dat, &payload)
	check(err)

	/* 	// Unmarshal fie with Decoder
	   	decoder := xml.NewDecoder(dat)
	   	for {
	   		// Read tokens from the XML document in a stream.
	   		t, _ := decoder.Token()
	   		if t == nil {
	   			break
	   		}
	   		// Inspect the type of the token just read.
	   		switch se := t.(type) {
	   		case xml.StartElement:
	   			// If we just read a StartElement token
	   			// ...and its name is "page"
	   			if se.Name.Local == "page" {
	   				var p Page
	   				// decode a whole chunk of following XML into the
	   				// variable p which is a Page (se above)
	   				decoder.DecodeElement(&p, &se)
	   				// Do some stuff with the page.
	   				p.Title = CanonicalizeTitle(p.Title)
	   				...
	   			}
	   	} */

	fmt.Println(payload.ProfileBody.DeviceIdentity.DeviceId)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
