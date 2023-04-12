package model

import "encoding/xml"

type RequestEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  struct{} `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body    struct {
		Add struct {
			A int `xml:"http://example.com/calculator a"`
			B int `xml:"http://example.com/calculator b"`
		} `xml:"http://example.com/calculator Add"`
	} `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type ResponseEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Header  struct{} `xml:"soap:Header"`
	Body    struct {
		AddResponse struct {
			Result int `xml:"tns:result"`
		} `xml:"tns:AddResponse"`
	} `xml:"soap:Body"`
}

type SoapFault struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Header  struct{} `xml:"soap:Header"`
	Body    struct {
		Fault struct {
			Faultcode   string `xml:"faultcode"`
			Faultstring string `xml:"faultstring"`
			Faultdetail struct {
				Fault struct {
					Code   string `xml:"tns:code"`
					Reason string `xml:"tns:reason"`
				} `xml:"tns:CalculatorFault"`
			} `xml:"detail"`
		} `xml:"soap:Fault"`
	} `xml:"soap:Body"`
}
