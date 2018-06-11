package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"strings"
	"encoding/xml"
	"io/ioutil"
	"github.com/tealeg/xlsx"
)

const (
	pathAndroidTH = "android/values"
	pathAndroidEN = "android/values-en"
	pathIosTH = "ios/th"
	pathIosEN = "ios/en"
	fileNameAndroid = "message.xml"
	fileNameIos = "Localizable.strings"
)

var (
	input = flag.String("f", "", "-f=ListOfValue.xlsm")
)

type Resources struct {
	XMLName 	xml.Name `xml:"resources"`
	Name []Msg 
}

type Msg struct {
	XMLName 	xml.Name `xml:"string"`
	ID string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func main(){
	flag.Parse()

	if *input == "" {
		fmt.Println("usage: go run main.go -f=LiveOfValue.xlsx ")
	}

	execForAndroid()
}

func execForAndroid(){
	en := Resources{}
	th := Resources{}
	iosEn := []string{}
	iosTh := []string{}

	execFileName := *input
	xlFile, err := xlsx.OpenFile(execFileName)

	if err != nil {
		log.Fatal(err)
	}

	for _, sheet := range xlFile.Sheets{
		if sheet.Name == "Sheet1" {
			body := sheet.Rows[1:]
			
			e := []Msg{}
			t := []Msg{}
			
			for _, row := range body {
				for i := len(row.Cells); i <= 3; i++ {
					row.Cells = append(row.Cells, &xlsx.Cell{})
				}
				
				e = append(e, Msg{ID:row.Cells[0].String(), Value: row.Cells[2].String()})
				t = append(t, Msg{ID:row.Cells[0].String(), Value: row.Cells[1].String()})
				iosEn = append(iosEn, fmt.Sprintf(`"%s"="%s";`, row.Cells[0].String(), row.Cells[2].String()))
				iosTh = append(iosTh, fmt.Sprintf(`"%s"="%s";`, row.Cells[0].String(), row.Cells[1].String()))
			}

			en.Name = e
			th.Name = t

			err := os.MkdirAll(pathAndroidTH, 0777);
			if err != nil {
				log.Fatal(err)
			}

			err = os.MkdirAll(pathAndroidEN, 0777);
			if err != nil {
				log.Fatal(err)
			}

			err = os.MkdirAll(pathIosTH, 0777);
			if err != nil {
				log.Fatal(err)
			}

			err = os.MkdirAll(pathIosEN, 0777);
			if err != nil {
				log.Fatal(err)
			}
			
			
			if xmlstring, err := xml.MarshalIndent(th, "", "    "); err == nil {
				xmlstring = []byte(xml.Header + string(xmlstring))
				fmt.Printf("%s\n",xmlstring)
				ioutil.WriteFile(pathAndroidTH + "/" + fileNameAndroid, []byte(xmlstring), 0777)
			}

			if xmlstring, err := xml.MarshalIndent(en, "", "    "); err == nil {
				xmlstring = []byte(xml.Header + string(xmlstring))
				fmt.Printf("%s\n",xmlstring)
				ioutil.WriteFile(pathAndroidEN + "/" + fileNameAndroid, []byte(xmlstring), 0777)
			}

			ioutil.WriteFile(pathIosEN + "/" + fileNameIos, []byte(strings.Join(iosEn, "\n")), 0777)
			ioutil.WriteFile(pathIosTH + "/" + fileNameIos, []byte(strings.Join(iosTh, "\n")), 0777)
		}
	}
}