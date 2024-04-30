package main

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	playfab "github.com/dgkanatsios/playfabsdk-go/sdk"
	client "github.com/dgkanatsios/playfabsdk-go/sdk/client"
)

func httpreq(url string) {
	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating GET request: %s", err)
	}

	// Add the X-Token header to the request
	req.Header.Set("X-Token", "7mQ6nz3umb7nPCcajUUsvNxlE9dZas37Scr930jaVtA=")

	// Perform the HTTP request
	client := &http.Client{}
	// get headers from req
	fmt.Println(req.Header)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing GET request: %s", err)
	}
	// get response header
	fmt.Println(response.Header)
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	fmt.Println(string(body))
}
func login(settings *playfab.Settings) (*client.LoginResultModel, error) {
	loginData := &client.LoginWithCustomIDRequestModel{
		CustomId:      "GettingStartedGuide",
		CreateAccount: true,
	}
	res, err := client.LoginWithCustomID(settings, loginData)
	return res, err
}

func loginWithGameCenter(settings *playfab.Settings) (*client.LoginResultModel, error) {
	loginData := &client.LoginWithGameCenterRequestModel{
		PlayerId:      "GameCenterPlayerID",
		CreateAccount: true,
	}
	res, err := client.LoginWithGameCenter(settings, loginData)
	return res, err
}

func DecodeSaveData(SaveData string) []byte {
	//SaveData := "H4s…"
	sDec, _ := base64.StdEncoding.DecodeString(SaveData)
	rData := strings.NewReader(string(sDec))
	gz, err := gzip.NewReader(rData)
	if err != nil {
		log.Fatal(err)
	}
	pdpData, _ := io.ReadAll(gz)
	fmt.Println(string(pdpData))
	return pdpData
}
func ParseJSONData(data []byte) map[string]interface{} {
	var jsonData map[string]interface{} // Or use a struct that matches your JSON structure
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	//return jsonData
	return jsonData

	// Now you have the parsed JSON data in the `jsonData` variable
	// You can further process it as needed
}

func main() {
	httpreq("https://google.com")
	settings := playfab.NewSettingsWithDefaultOptions("8b9b7")

	res, err := loginWithGameCenter(settings)
	if err != nil {
		fmt.Printf("Login Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Login OK, SessionTicket: %s\n", res.SessionTicket)

	sessionTicket := res.SessionTicket
	accountInfoRequest := &client.GetAccountInfoRequestModel{}

	accountinfo, err := client.GetAccountInfo(settings, accountInfoRequest, sessionTicket)

	if err != nil {
		fmt.Printf("Get account Error: %s \n", err.Error())
		return
	}

	//fmt.Printf("Get account OK: %#v\n", accountinfo.AccountInfo)
	// parse account info with getaccountinforesultmodel
	fmt.Printf("Get account OK: %#v\n", accountinfo.AccountInfo)

	fmt.Printf("Get gamecenter OK: %#v\n", accountinfo.AccountInfo.GameCenterInfo)
	fmt.Printf("Get UserTitleInfo OK: %#v\n", accountinfo.AccountInfo.TitleInfo)
	fmt.Printf("Get private info OK: %#v\n", accountinfo.AccountInfo.PrivateInfo)
	// get title data
	titleDataRequest := &client.GetTitleDataRequestModel{}
	titleData, err := client.GetTitleData(settings, titleDataRequest, sessionTicket)
	if err != nil {
		fmt.Printf("Get title data Error: %s\n", err.Error())
		return
	}
	//fmt.Printf("Get title data OK: %#v\n", titleData.Data)
	assetBundle := titleData.Data["AssetBundleInfos_1.34.0"]
	abData := titleData.Data["AbData_1.34.0"]
	flags := titleData.Data["Flags"]
	// if abData is null, get previous version
	if abData == "" {
		abData = titleData.Data["AbData_1.33.0"]
	}
	fmt.Printf("Get AssetBundles_1.34.0 OK: %#v\n", assetBundle)
	fmt.Printf("Get AbData_1.34.0 OK: %#v\n", abData)
	fmt.Printf("Get Flags OK: %#v\n", flags)
	// get object LteSchedule_1.34.0 from title data
	lteSchedule := titleData.Data["LteSchedule_1.34.0"]
	//fmt.Printf("Get LteSchedule_1.34.0 OK: %#v\n", lteSchedule)
	// decode lteschedule
	decodedData := DecodeSaveData(lteSchedule)
	jsonObject := ParseJSONData([]byte(decodedData))
	// get keys from jsonobject[LteDatas]
	events := jsonObject["LteDatas"].([]interface{})
	// check if events is not nil before accessing its elements
	// iterate over events
	for _, event := range events {
		// type assert event to map[string]interface{}
		eventMap := event.(map[string]interface{})
		// get value of id, gamedataid, startdatetimeutc and enddatetimeutc from eventMap
		id := eventMap["Id"]
		gamedataid := eventMap["GameDataId"]
		startdatetimeutc := eventMap["StartDateTimeUtc"]
		enddatetimeutc := eventMap["EndDateTimeUtc"]
		fmt.Printf("Id: %s, GameDataId: %s, StartDateTimeUtc: %s, EndDateTimeUtc: %s\n", id, gamedataid, startdatetimeutc, enddatetimeutc)

	}

}
