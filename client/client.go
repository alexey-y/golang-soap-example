package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Header  struct{} `xml:"soap:Header"`
	Body    struct {
		Add struct {
			A int `xml:"tns:a"`
			B int `xml:"tns:b"`
		} `xml:"tns:Add"`
	} `xml:"soap:Body"`
}

type ResponseEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  struct{} `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body    struct {
		AddResponse struct {
			Result int `xml:"http://example.com/calculator result"`
		} `xml:"http://example.com/calculator AddResponse"`
	} `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type SoapFault struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  struct{} `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body    struct {
		Fault struct {
			Code   string `xml:"faultcode"`
			String string `xml:"faultstring"`
			Detail struct {
				CalculatorFault struct {
					Code   string `xml:"http://example.com/calculator code"`
					Reason string `xml:"http://example.com/calculator reason"`
				} `xml:"http://example.com/calculator CalculatorFault"`
			} `xml:"detail"`
		} `xml:"http://www.w3.org/2003/05/soap-envelope Fault"`
	} `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type SoapError struct {
	Code   string
	Reason string
}

func (se SoapError) Error() string {
	return fmt.Sprintf("SOAP error with code %s: %s", se.Code, se.Reason)
}

func CallAddService(a, b int) (int, error) {
	reqEnv := RequestEnvelope{}
	reqEnv.XMLName.Space = "http://www.w3.org/2003/05/soap-envelope"
	reqEnv.XMLName.Local = "Envelope"
	reqEnv.Header = struct{}{}
	reqEnv.Body.Add.A = a
	reqEnv.Body.Add.B = b

	reqXML, err := xml.MarshalIndent(reqEnv, "", "  ")
	if err != nil {
		return 0, err
	}

	customXMLHeader := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:tns="http://example.com/calculator">`

	reqXMLString := strings.Replace(string(reqXML), "<soap:Envelope>", customXMLHeader, 1)

	req, err := http.NewRequest("POST", "http://localhost:8080/soap", bytes.NewReader([]byte(reqXMLString)))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://example.com/calculator/Add")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respXML, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	hasFault := strings.Contains(string((respXML)), "soap:F")

	if hasFault {
		respFault := SoapFault{}
		err = xml.Unmarshal(respXML, &respFault)
		if err != nil {
			return 0, err
		}
		return -1, SoapError{Code: respFault.Body.Fault.Detail.CalculatorFault.Code, Reason: respFault.Body.Fault.Detail.CalculatorFault.Reason}
	}

	respEnv := ResponseEnvelope{}
	err = xml.Unmarshal(respXML, &respEnv)
	if err != nil {
		return 0, err
	}

	return respEnv.Body.AddResponse.Result, nil
}
