package server_test

import (
	"fmt"
	"net/http"
	"soap-example/client"
	"soap-example/server"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var once sync.Once

func startServer() {
	once.Do(func() {
		http.HandleFunc("/soap", func(w http.ResponseWriter, r *http.Request) {
			server.SoapHandler(w, r, &server.CalculatorService{})
		})
		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				fmt.Println("Failed to start server:", err)
			}
		}()
	})
}

func TestAddition(t *testing.T) {

	startServer()
	Convey("Given a running SOAP service", t, func() {
		Convey("When a client calls the Add method with positive numbers", func() {
			a, b := 12, 18
			result := 30

			result, err := client.CallAddService(a, b)

			Convey("The server should returns the sum", func() {
				So(err, ShouldBeNil)
				So(result, ShouldEqual, result)
			})
		})
	})
}

func TestFault(t *testing.T) {

	startServer()
	Convey("Given a running SOAP service", t, func() {

		Convey("When a client calls the Add method with negative numbers", func() {
			a, b := -1, 18

			_, err := client.CallAddService(a, b)

			Convey("The server should return SoapFault", func() {
				customErr, isSoapError := err.(client.SoapError)
				So(isSoapError, ShouldBeTrue)
				So(customErr.Code == "C101", ShouldBeTrue)
			})
		})
	})
}
