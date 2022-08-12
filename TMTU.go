package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	geojson "github.com/paulmach/go.geojson"
)

type ResponseTMTU struct {
	Status   string `json:"status"`
	Messages string `json:"messages"`
	Data     []struct {
		WpointName    string `json:"WpointName"`
		WpointNameEng string `json:"WpointNameEng"`
		WPointNo      string `json:"WPointNo"`
		Longitude     string `json:"Longitude"`
		Latitude      string `json:"Latitude"`
		GroupType     string `json:"group_type"`
	} `json:"data"`
}

type XYZData struct {
    Tags []string `json:"tags"`
}

func main() {

	respTMTU, err := http.Get("http://tmtitsapi.locationtracker.com/api/getWayPoints")		//GET request to TMTU for ticker prices and withdrawal/deposit enabled lists
	if err != nil {
		log.Fatal(err)
	}
	defer respTMTU.Body.Close()
	bodyTMTU, err := ioutil.ReadAll(respTMTU.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultTMTU ResponseTMTU
	if err := json.Unmarshal(bodyTMTU, &resultTMTU); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}
	fc := geojson.NewFeatureCollection()
	//fmt.Println(resultTMTU.Data[0].WpointName)

	
	for i := 0; i < 721; i++ {
		sresultTMTULatitude, err := strconv.ParseFloat(resultTMTU.Data[i].Latitude, 64)
		if err != nil {
			fmt.Println("wrong here1")
		}
		sresultTMTULongitude, err := strconv.ParseFloat(resultTMTU.Data[i].Longitude, 64)
		if err != nil {
			fmt.Println("wrong here1")
		}
    	feature := geojson.NewPointFeature([]float64{sresultTMTULongitude, sresultTMTULatitude})
   	 	feature.SetProperty("name", resultTMTU.Data[i].WpointName)
		feature.SetProperty("ref", resultTMTU.Data[i].WPointNo)
		feature.SetProperty("highway", "bus_stop")
		feature.SetProperty("operator", "Thane Municipal Transport")
		feature.SetProperty("public_transport", "platform")
  		fc.AddFeature(feature)
	}
	
	
	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//fmt.Printf("%s", string(rawJSON))

	err = ioutil.WriteFile("userfile.json", rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

}