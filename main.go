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

	funcMap := template.FuncMap{
		// funkcja licząca czy po tej iteracji zapisywać drukowanie już na następnej stronie
		"check": func(i int) bool {
			return (i % 3) == 2
		},
	}
	dest := os.Args[2]
	if dest == "" {
		dest = "result.html"
	}
	file, _ := os.Create(dest)
	/*templatka naszego HTML*/
	t, err := template.New("res").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
    <style>
        .item { width: 30%; float: left; margin: 1%; }
        .page-break { clear: both; page-break-after: always; }
    </style>
</head>
<body style="font-size: 85%;">
	{{$header := .Message.DD801B.Header}}
	{{$body := .Message.DD801B.Body}}
    {{range $index, $item := .Message.DD801B.Body.EDDContainer.ConsigneeTraders}}
        
	<table border="1" cellpadding="5" cellspacing="0" width="100%">
    <tr>
        <td>
            Dokument e-DD
        </td>
    </tr>
    <tr>
        <td align="right">
            <small>Identyfikator komunikatu:</small> {{ $header.MessageIdentifier}}, utworzonego: {{$header.DateOfPreparation}} {{$header.TimeOfPreparation}}<br>
            <small>Numer LRN przemieszczenia:</small> {{ $body.EDDContainer.EDD.LocalReferenceNumber}}<br>
            <small>Numer ARC przemieszczenia:</small> {{ $body.EDDContainer.EDD.DeliveryDocumentReference.DeliveryDocumentAdministrativeReferenceCode}}<br>
            <small>Data i czas pierwszej walidacji projektu e-DD:</small> {{$body.EDDContainer.EDD.DateAndTimeOfValidationOfEDD}}
        </td>
    </tr>
    <tr>
        <td>
            <table width="100%" cellpadding="5" cellspacing="0"  style="border-collapse: collapse;">
                <tr>
                    <td width="50%" valign="top">

                        <strong>Podmiot odbierający</strong><br>

                    </td>
                    <td width="50%" valign="top" style="border-left: 1px solid black;">
                        <strong>Podmiot miejsce dostawy</strong><br>
                    </td>
                </tr>
            </table>
            <table width="100%" cellpadding="5" cellspacing="0"  style="border-collapse: collapse;">
                <tr>
                    <td width="50%" valign="top">
                        {{ $item.ConsigneeTrader.TraderName}}<br>
                        {{ $item.ConsigneeTrader.StreetName }} {{$item.ConsigneeTrader.StreetNumber}},<br>
                        {{ $item.ConsigneeTrader.Postcode }} {{ $item.ConsigneeTrader.City }}<br>
                        Identyfikator podmiotu: {{ $item.ConsigneeTrader.TraderId.ExciseNumber }}
                    </td>
                    <td width="50%" valign="top" style="border-left: 1px solid black;">
                    </td>
                </tr>
            </table>
        </td>
    </tr>
    <tr>
        <td>
            Odebrano litrów: ....................<br>
            Czytelny podpis (imię i nazwisko): ................................................................................
        </td>
    </tr>
    <tr>
        <td>
            Oznaczenie jednostek transp. <b>DSW PY11</b>  Kod wyrobu akcyzowego: <b>{{$body.EDDContainer.BodyEDD.ExciseProductCode}}</b> Kod CN: <b>{{$body.EDDContainer.BodyEDD.CnCode}}</b>
            Ilość: <b>{{$body.EDDContainer.BodyEDD.Quantity}}</b> Masa brutto: <b>{{$body.EDDContainer.BodyEDD.GrossWeight}}</b> Masa netto: <b>{{$body.EDDContainer.BodyEDD.NetWeight}}</b>Opis handlowy [pl]: <b>{{$body.EDDContainer.BodyEDD.CommercialDescription.Value}}</b>
        </td>
    </tr>
</table>
        {{if (check $index)}}
            <div class="page-break"></div>
        {{end}}
			
    {{end}}
</body>
</html>`) // do jednej templatki nie opłaca się tworzyć osobnego folderu z templatkami
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
		fmt.Printf("za mało argumentów\n")
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
	sanitizeEdd(&edd)
	res := makeHtml(edd)
	fmt.Printf("udało się stworzyć dane do wydruku pod %v\n", res)
}
