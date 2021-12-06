package main

import (
	"go.uber.org/zap"
	"time"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"net/http"
)

// DeviceIdentity : IODevice>ProfileBody>DeviceIdentity
type DeviceIdentity struct {
	VendorId int `xml:"vendorId,attr"`
	VendorName string `xml:"vendorName,attr"`
	DeviceId int `xml:"deviceId,attr"` // id of type of device, given by device vendor
}



// ProcessData : IODevice>ProfileBody>DeviceFunction>ProcessDataCollection>ProcessData
type ProcessData struct {
	
}

type IoddInformationOfOneVendor struct {
	vendorId string
	devicesOfVendor []Device
}

type Device struct {
	deviceTypeID int

}

//https://jm33.me/parsing-large-xml-with-go.html


func processIoddFile(file){
	decoder := xml.NewDecoder(r io.Reader)
	for {
		t, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
			break
			}
			// handle error somehow
			return fmt.Errorf("decoding token: %v", err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Space == "foo" && t.Name.Local == "bar" {
				var b bar
				if err := decoder.DecodeElement(&b, &t); err != nil {
					// handle error somehow
					return fmt.Errorf("decoding element %q: %v", t.Name.Local, err)
				}
				// do something with b
			}
		}
	}

	func getStructMember(parser *xml.Decoder) (member Struct) {
		var token xml.Token
		token, _ = parser.Token()
	
		member = Struct{}
	
		for {
			switch t := token.(type) {
			case xml.StartElement:
				if t.Name.Local == "name" {
					member["name"], _ = getElementValue(parser)
				}
	
				if t.Name.Local == "value" {
					member["value"], _ = getValue(parser)
				}
			case xml.EndElement:
				if t.Name.Local == "member" {
					return member
				}
			}
	
			token, _ = parser.Token()
		}
	
		return
	}


	
	<?xml version="1.0" encoding="ISO-8859-1" ?>
	<FileRetriever>
	  <FileList>
		  <File name="Name1" />
		  <File name="Name2" />
	  </FileList>
	</FileRetriever>
	type fileRetriever struct {
		Files []file `xml:"FileList>File"`
	}
	
	type file struct {
		Name string `xml:"name,attr"`
	}
	
	func Main(){
		retrieve()
	}
	
	func retrieve()(retriever *fileRetriever){
		req := ... //set up http.NewRequest()
		client := &http.Client{}
		rsp, err := client.Do(req)
	
		if err != nil {
			log.Fatal(err)
		}
	
		defer rsp.Body.Close()
	
		decoder := xml.NewDecoder(rsp.Body)
		decoder.CharsetReader = charset.NewReaderLabel
	
		retriever = &fileRetriever{}
	
		err = decoder.Decode(&retriever)
	
		if err != nil {
			fmt.Println(err)
		}
	
		return retriever, xTidx
	}
}