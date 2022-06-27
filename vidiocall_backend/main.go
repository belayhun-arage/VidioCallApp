package main

import (
	controller "github.com/belayhun-arage/VidioCallApp/vidiocall_backend/controller"
	"fmt"
	"log"
	"net/http"
)

func main(){
	// Handling routes
	// RTC token from RTC num uid
	var c controller.IRTCTokenController=controller.NewRTCTokenController();
	
	http.HandleFunc("/fetch_rtc_token", c.NewUserRtcTokenHandler)
	fmt.Printf("Starting server at port 8082\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
 
