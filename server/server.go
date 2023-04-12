package server

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"soap-example/server/model"
	"strings"
)

type CalculatorService struct{}

func (s *CalculatorService) Add(a int, b int) int {
	return a + b
}

func SoapHandler(w http.ResponseWriter, r *http.Request, service *CalculatorService) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	request := model.RequestEnvelope{}
	err = xml.Unmarshal(bodyBytes, &request)
	if err != nil {
		http.Error(w, "Failed to parse request XML", http.StatusBadRequest)
		return
	}

	a := request.Body.Add.A
	b := request.Body.Add.B

	if a < 0 || b < 0 {
		sendSoapFault(w, "C101", "Negative a or b")
		return
	}

	result := service.Add(a, b)

	response := model.ResponseEnvelope{}
	response.Body.AddResponse.Result = result

	responseBytes, err := xml.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate response XML", http.StatusInternalServerError)
		return
	}
	customXMLHeader := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:tns="http://example.com/calculator">`

	responseString := strings.Replace(string(responseBytes), "<soap:Envelope>", customXMLHeader, 1)

	w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
	w.Write([]byte(responseString))
}

func sendSoapFault(w http.ResponseWriter, faultCode, faultReason string) {
	soapFault := model.SoapFault{}
	soapFault.Body.Fault.Faultcode = "soap:Client"
	soapFault.Body.Fault.Faultstring = "Client error"
	soapFault.Body.Fault.Faultdetail.Fault.Code = faultCode
	soapFault.Body.Fault.Faultdetail.Fault.Reason = faultReason

	responseBytes, err := xml.Marshal(soapFault)
	if err != nil {
		http.Error(w, "Failed to generate response XML", http.StatusInternalServerError)
		return
	}
	customXMLHeader := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:tns="http://example.com/calculator">`

	responseString := strings.Replace(string(responseBytes), "<soap:Envelope>", customXMLHeader, 1)

	w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(responseString))
}
