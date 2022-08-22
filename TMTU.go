package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"net/url"
	geojson "github.com/paulmach/go.geojson"
)

type ResponseWaypoints struct {
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

type ResponseRouteMaster struct {
	Status   string `json:"status"`
	Messages string `json:"messages"`
	Data     []struct {
		RouteNo        string `json:"RouteNo"`
		RouteName      string `json:"RouteName"`
		RouteNum       string `json:"RouteNum"`
		RouteDirection string `json:"RouteDirection"`
	} `json:"data"`
}

type ResponseRouteNo struct {
	Status   string `json:"status"`
	Messages string `json:"messages"`
	Data     []struct {
		RouteNo                 int         `json:"RouteNo"`
		RouteName               string      `json:"RouteName"`
		RouteNum                string      `json:"RouteNum"`
		RouteDirection          string      `json:"RouteDirection"`
		RouteStage              interface{} `json:"RouteStage"`
		TotalCalculatedDistance string      `json:"total_calculated_distance"`
		RouteDetails            []struct {
			RDNo       int    `json:"RDNo"`
			RouteNo    string `json:"RouteNo"`
			WPointNo   string `json:"WPointNo"`
			SequenceNo string `json:"SequenceNo"`
			Waypoints  struct {
				WPointNo     string        `json:"WPointNo"`
				WpointName   string        `json:"WpointName"`
				Longitude    string        `json:"Longitude"`
				Latitude     string        `json:"Latitude"`
				InsertedDate string        `json:"InsertedDate"`
				GroupType    string        `json:"group_type"`
				InRouteNo    string        `json:"in_route_no"`
				IsSuspected  string        `json:"is_suspected"`
				Allvehicle   []interface{} `json:"allvehicle"`
			} `json:"waypoints"`
		} `json:"route_details"`
	} `json:"data"`
	AllRouteVehicles []interface{} `json:"all_route_vehicles"`
}

func main() {
	waypoints()
	routes()
}

func waypoints() {
	respWaypoints, err := http.Get("http://tmtitsapi.locationtracker.com/api/getWayPoints") //GET request to TMTU for waypoints data
	if err != nil {
		log.Fatal(err)
	}
	defer respWaypoints.Body.Close()
	bodyWaypoints, err := io.ReadAll(respWaypoints.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultWaypoints ResponseWaypoints
	if err := json.Unmarshal(bodyWaypoints, &resultWaypoints); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}
	waypoints := geojson.NewFeatureCollection()
	//fmt.Println(resultWaypoints.Data[0].WpointName)

	//Convert the data to geojson format for JOSM
	for i := 0; i < 721; i++ {
		sresultWaypointsLatitude, err := strconv.ParseFloat(resultWaypoints.Data[i].Latitude, 64)
		if err != nil {
			fmt.Println("wrong here1")
		}
		sresultWaypointsLongitude, err := strconv.ParseFloat(resultWaypoints.Data[i].Longitude, 64)
		if err != nil {
			fmt.Println("wrong here1")
		}
		feature := geojson.NewPointFeature([]float64{sresultWaypointsLongitude, sresultWaypointsLatitude})
		feature.SetProperty("name", resultWaypoints.Data[i].WpointName)
		feature.SetProperty("ref", resultWaypoints.Data[i].WPointNo)
		feature.SetProperty("highway", "bus_stop")
		feature.SetProperty("operator", "Thane Municipal Transport")
		feature.SetProperty("public_transport", "platform")
		waypoints.AddFeature(feature)
	}

	rawJSON, err := waypoints.MarshalJSON()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//Saves the geojson file for bus stops to the current directory
	//fmt.Printf("%s", string(rawJSON))
	err = os.WriteFile("TMTStops.json", rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func routes() {

	respRoutes, err := http.Get("http://tmtitsapi.locationtracker.com/api/getRouteMaster") //GET request to TMTU for routes data
	if err != nil {
		log.Fatal(err)
	}
	defer respRoutes.Body.Close()
	bodyRoutes, err := io.ReadAll(respRoutes.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultRoutes ResponseRouteMaster
	if err := json.Unmarshal(bodyRoutes, &resultRoutes); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	for i := 0; i < 158; i++ {

		data := url.Values{
			"RouteNo":       {resultRoutes.Data[i].RouteNo},
		}
		respRouteNo, err := http.PostForm("http://tmtitsapi.locationtracker.com/api/getRouteDetailsNew", data) //GET request to TMTU for routes data
		if err != nil {
			log.Fatal(err)
		}
		defer respRouteNo.Body.Close()
		bodyRouteNo, err := io.ReadAll(respRouteNo.Body) // response body is []byte
		if err != nil {
			fmt.Println("wrong here")
		}
		var resultRouteNo ResponseRouteNo
		if err := json.Unmarshal(bodyRouteNo, &resultRouteNo); err != nil { // Parse []byte to the go struct pointer
			fmt.Println(err)
		}
		fmt.Print(resultRouteNo)
		routes := geojson.NewFeatureCollection()
		for j := 0; j < len(resultRouteNo.Data); j++ {
			sresultRouteNoLatitude, err := strconv.ParseFloat(resultRouteNo.Data[0].RouteDetails[j].Waypoints.Latitude, 64)
			if err != nil {
				fmt.Println("wrong here1")
			}
			sresultRouteNoLongitude, err := strconv.ParseFloat(resultRouteNo.Data[0].RouteDetails[j].Waypoints.Longitude, 64)
			if err != nil {
				fmt.Println("wrong here1")
			}

			feature := geojson.NewLineStringFeature([][]float64{sresultRouteNoLatitude, sresultRouteNoLongitude},)
			feature.SetProperty("name", resultRouteNo.Data[i].RouteDetails[j].Waypoints.WpointName)
			routes.AddFeature(feature)
		}
		
	}
}
