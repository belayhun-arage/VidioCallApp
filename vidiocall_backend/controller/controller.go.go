package controller

import(
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	"github.com/belayhun-arage/VidioCallApp/vidiocall_backend/entity"
	"net/http"
	"fmt"
	"log"
	"time"
	"errors"
	"strconv"
	"os"
	"encoding/json"
)

var rtc_token string
var int_uid uint32
var channel_name string

var role_num uint32
var role rtctokenbuilder.Role

type IRTCTokenController interface{
	NewUserRtcTokenHandler(w http.ResponseWriter, r *http.Request)
}

type RTCTokenController struct{}

func NewRTCTokenController() IRTCTokenController {
	return &RTCTokenController{}
}


func generateRtcToken(int_uid uint32, channelName string, role rtctokenbuilder.Role){

	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("APP_CERTIFICATE")
	// Number of seconds after which the token expires.
	// For demonstration purposes the expiry time is set to 40 seconds. This shows you the automatic token renew actions of the client.
	expireTimeInSeconds := uint32(40)
	// Get current timestamp.
	currentTimestamp := uint32(time.Now().UTC().Unix())
	// Timestamp when the token expires.
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, int_uid, role, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with uid: %s\n", result)
		fmt.Printf("uid is %d\n", int_uid )
		fmt.Printf("ChannelName is %s\n", channelName)
		fmt.Printf("Role is %d\n", role)
	}
	rtc_token = result
}

func(*RTCTokenController) NewUserRtcTokenHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
	w.Header().Set("Access-Control-Allow-Headers", "*");

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
		return
	}


	var t_int entity.RTCTocken
	var unmarshalErr *json.UnmarshalTypeError
	int_decoder := json.NewDecoder(r.Body)
	int_err := int_decoder.Decode(&t_int)
	if (int_err == nil) {

				int_uid = t_int.User_Id
				channel_name = t_int.Channel_Name
				role_num = t_int.Role
				switch role_num {
				case 0:
					// DEPRECATED. RoleAttendee has the same privileges as RolePublisher.
					role = rtctokenbuilder.RoleAttendee
				case 1:
					role = rtctokenbuilder.RolePublisher
				case 2:
					role = rtctokenbuilder.RoleSubscriber
				case 101:
					// DEPRECATED. RoleAdmin has the same privileges as RolePublisher.
					role = rtctokenbuilder.RoleAdmin
				}
	}
	if (int_err != nil) {

		if errors.As(int_err, &unmarshalErr){
				errorResponse(w, "Bad request. Wrong type provided for field " + unmarshalErr.Value  + unmarshalErr.Field + unmarshalErr.Struct, http.StatusBadRequest)
				} else {
				errorResponse(w, "Bad request.", http.StatusBadRequest)
			}
		return
	}

	generateRtcToken(int_uid, channel_name, role)
	errorResponse(w, rtc_token, http.StatusOK)
	log.Println(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}