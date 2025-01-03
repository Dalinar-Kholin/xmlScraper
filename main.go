package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
)

type EDDToTrader struct {
	Message Message `xml:"Message"`
}

type Message struct {
	DD801B DD801B `xml:"http://www.mf.gov.pl/schematy/EDD/DD801B/2022/07/ DD801B"`
}

type DD801B struct {
	Header Header `xml:"Header"`
	Body   Body   `xml:"Body"`
}

type Header struct {
	DateOfPreparation string `xml:"http://www.mf.gov.pl/schematy/EDD/DDHEADERTYPE/2022/07/ DateOfPreparation"`
	TimeOfPreparation string `xml:"http://www.mf.gov.pl/schematy/EDD/DDHEADERTYPE/2022/07/ TimeOfPreparation"`
	MessageIdentifier string `xml:"http://www.mf.gov.pl/schematy/EDD/DDHEADERTYPE/2022/07/ MessageIdentifier"`
}

type Body struct {
	EDDContainer EDDContainer `xml:"EDDContainer"`
}

type EDDContainer struct {
	EDD              EDD                `xml:"EDD"`
	ConsigneeTraders []ConsigneeTraders `xml:"ConsigneeTraders"`
	TransportMode    TransportMode      `xml:"TransportMode"`
	BodyEDD          BodyEDD            `xml:"BodyEDD"`
}

type EDD struct {
	DeliveryDocumentReference    DeliveryDocumentReference `xml:"DeliveryDocumentReference"`
	LocalReferenceNumber         string                    `xml:"LocalReferenceNumber"`
	DateAndTimeOfValidationOfEDD string                    `xml:"DateAndTimeOfValidationOfEDD"`
}

type DeliveryDocumentReference struct {
	XMLName                                     xml.Name `xml:"DeliveryDocumentReference"`
	DeliveryDocumentAdministrativeReferenceCode string   `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ DeliveryDocumentAdministrativeReferenceCode"`
}

type TraderId struct {
	ExciseNumber string `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ ExciseNumber,omitempty"`
	TaxNumber    string `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ TaxNumber,omitempty"`
	PersonalId   string `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ PersonalId,omitempty"`
}

type CustomsOffice struct {
	ReferenceNumber string `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ ReferenceNumber,omitempty"`
}
type DeliveryPlaceTrader struct {
	TraderId     TraderId `xml:"Traderid"`
	TraderName   string   `xml:"TraderName"`
	StreetName   string   `xml:"StreetName"`
	StreetNumber string   `xml:"StreetNumber"`
	Postcode     string   `xml:"Postcode"`
	City         string   `xml:"City"`
}

type ConsigneeTraders struct {
	ConsigneeTrader            ConsigneeTrader     `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ ConsigneeTrader"`
	DeliveryPlaceCustomsOffice CustomsOffice       `xml:"DeliveryPlaceCustomsOffice,omitempty"`
	DeliveryPlaceTrader        DeliveryPlaceTrader `xml:"DeliveryPlaceTrader,omitempty"`
}

type ConsigneeTrader struct {
	DeliveryTraderType string   `xml:"deliveryTraderType,attr"`
	Language           string   `xml:"language,attr"`
	TraderId           TraderId `xml:"Traderid"`
	TraderName         string   `xml:"TraderName"`
	StreetName         string   `xml:"StreetName"`
	StreetNumber       string   `xml:"StreetNumber"`
	Postcode           string   `xml:"Postcode"`
	City               string   `xml:"City"`
}

type TransportMode struct {
	TransportModeCode string `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ TransportModeCode"`
}
type BodyEDD struct {
	ExciseProductCode     string                `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ ExciseProductCode"`
	CnCode                string                `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ CnCode"`
	Quantity              string                `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ Quantity"`
	GrossWeight           string                `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ GrossWeight"`
	NetWeight             string                `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ NetWeight"`
	CommercialDescription CommercialDescription `xml:"http://www.mf.gov.pl/schematy/EDD/DDTYPES/2022/07/ CommercialDescription"`
}

type CommercialDescription struct {
	Language string `xml:"language,attr"`
	Value    string `xml:",chardata"`
}

func makeHtml(data EDDToTrader) string {

	for _, x := range data.Message.DD801B.Body.EDDContainer.ConsigneeTraders {
		fmt.Printf("identyfikator:= %v\n", x.ConsigneeTrader.TraderId.PersonalId)
	}
	funcMap := template.FuncMap{
		// funkcja licząca czy po tej iteracji zapisywać drukowanie już na następnej stronie
		"check": func(i int) bool {
			return (i % 3) == 2
		},
		"getId": func(d TraderId) (res string) { // brzydkie w chuj ale cóż zrobię
			if res = d.TaxNumber; res != "" {
				return
			}
			if res = d.ExciseNumber; res != "" {
				return
			}
			if res = d.PersonalId; res != "" {
				return
			}
			return "dad XML or me stupido"
		},
	}
	dest := os.Args[2]
	if dest == "" {
		dest = "result.html"
	}
	file, _ := os.Create(dest)
	/*templatka naszego HTML*/
	htmlTemplate, err := os.ReadFile("./template.gohtml")
	if err != nil {
		panic(err)
	}
	t, err := template.New("res").Funcs(funcMap).Parse(string(htmlTemplate)) // do jednej templatki nie opłaca się tworzyć osobnego folderu z templatkami
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(file, "res", data)
	if err != nil {
		panic(err)
	}
	return dest
}

func sanitizeEdd(edd *EDDToTrader) {
	for i, x := range edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders {
		/*jeżeli streetName jest puste oznacza to że istnieje inny adres pod który musi być dostawa, więc podmieniamy dane, ponieważ w teplatce czytamy z podstawowych*/
		if x.DeliveryPlaceTrader.StreetName != "" {
			/*ICOM robione poprzez direct odniesienia, ponieważ nie da sie łatwo iterować po pointerach */
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.TraderId.ExciseNumber = x.DeliveryPlaceTrader.TraderId.TaxNumber
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.TraderName = x.DeliveryPlaceTrader.TraderName
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.StreetName = x.DeliveryPlaceTrader.StreetName
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.StreetNumber = x.DeliveryPlaceTrader.StreetNumber
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.Postcode = x.DeliveryPlaceTrader.Postcode
			edd.Message.DD801B.Body.EDDContainer.ConsigneeTraders[i].ConsigneeTrader.City = x.DeliveryPlaceTrader.City
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	var edd EDDToTrader
	if err := xml.Unmarshal(file, &edd); err != nil {
		panic(err)
	}
	//sanitizeEdd(&edd)
	res := makeHtml(edd)
	fmt.Printf("udało się stworzyć dane do wydruku pod %v\n", res)
}
