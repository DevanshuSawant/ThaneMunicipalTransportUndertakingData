package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/option"
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

type ResponseBusLocations struct {
	Status   string `json:"status"`
	Messages string `json:"messages"`
	Data     []struct {
		IdxTrackidPk     int         `json:"idx_Trackid_pk"`
		VehID            string      `json:"VehId"`
		CmpID            string      `json:"CmpId"`
		LastTrackdt      string      `json:"LastTrackdt"`
		NCSent           string      `json:"NCSent"`
		CSent            string      `json:"CSent"`
		PrevTrackDt      string      `json:"PrevTrackDt"`
		LastNCSentDate   string      `json:"LastNCSentDate"`
		Longitude        string      `json:"Longitude"`
		Latitude         string      `json:"Latitude"`
		City             interface{} `json:"City"`
		Speed            string      `json:"Speed"`
		ImagePath        interface{} `json:"ImagePath"`
		AC               string      `json:"AC"`
		Ignition         string      `json:"Ignition"`
		AUX1             string      `json:"AUX1"`
		DI4              string      `json:"DI4"`
		Fuel             string      `json:"Fuel"`
		Temparature      string      `json:"Temparature"`
		WPointNo         interface{} `json:"WPointNo"`
		Odometer         string      `json:"Odometer"`
		Distance         string      `json:"Distance"`
		ETATime          string      `json:"ETATime"`
		ETARoute         string      `json:"ETARoute"`
		ETAOldTime       string      `json:"ETAOldTime"`
		Routeflag        string      `json:"routeflag"`
		ETARouteName     string      `json:"ETARouteName"`
		DirectionFrom    string      `json:"DirectionFrom"`
		DirectionTo      string      `json:"DirectionTo"`
		DispatchDateTime string      `json:"DispatchDateTime"`
		ETATime1         string      `json:"ETATime1"`
		ETAOldTime1      string      `json:"ETAOldTime1"`
		Routeflag1       interface{} `json:"routeflag1"`
		RouteNo          string      `json:"RouteNo"`
		WaybillNo        string      `json:"WaybillNo"`
		Lastwaypointid   interface{} `json:"lastwaypointid"`
		Token            string      `json:"token"`
		Avgspeed         string      `json:"avgspeed"`
		LatLong          string      `json:"LatLong"`
		GetVehicle       struct {
			Vehid string `json:"vehid"`
			VehNo string `json:"VehNo"`
		} `json:"get_vehicle"`
	} `json:"data"`
}

// BusLocation is a json-serializable type.
type BusLocation struct {
	Vehid         string      `json:"veh_id,omitempty"`
	RouteNo       string      `json:"route_no,omitempty"`
	Longitude     string      `json:"longitude,omitempty"`
	Latitude      string      `json:"latitude,omitempty"`
	From          string      `json:"from,omitempty"`
	To            string      `json:"to,omitempty"`
	Waypoint      interface{} `json:"waypoint,omitempty"`
	LastWaypoint  interface{} `json:"last_waypoint,omitempty"`
	LastTrackTime string      `json:"last_track_time,omitempty"`
	DispatchTime  string      `json:"dispatch_time,omitempty"`
	VehNo         string      `json:"veh_no,omitempty"`
	CheckedAt     string      `json:"checked_at,omitempty"`
}

func main() {
	err := os.MkdirAll("output", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	waypoints()
	routes()

	for {
		fmt.Printf("Started At:%s\n", time.Now())
		buslocations()
	}

}

func waypoints() {
	respWaypoints, err := http.Get("http://tmtitsapi.locationtracker.com/api/getWayPoints") //GET request to TMTU for waypoints data
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(respWaypoints.Body)
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
		waypoints.AddFeature(feature) //
	}

	rawJSON, err := waypoints.MarshalJSON()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//Saves the geojson file for bus stops to the current directory
	//fmt.Printf("%s", string(rawJSON))
	err = os.WriteFile("output/TMTStopsDirect.json", rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func routes() {
	var stops = make(map[int]string)
	var ref []int
	respRoutes, err := http.Get("http://tmtitsapi.locationtracker.com/api/getRouteMaster") //GET request to TMTU for routes data
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(respRoutes.Body)
	bodyRoutes, err := io.ReadAll(respRoutes.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultRoutes ResponseRouteMaster
	if err := json.Unmarshal(bodyRoutes, &resultRoutes); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}
	waypoints := geojson.NewFeatureCollection()
	for i := 0; i < len(resultRoutes.Data); i++ {

		data := url.Values{
			"RouteNo": {resultRoutes.Data[i].RouteNo},
		}
		time.Sleep(2 * time.Second)
		fmt.Println("Restarting...")
		fmt.Println(i)
		respRouteNo, err := http.PostForm("http://tmtitsapi.locationtracker.com/api/getRouteDetailsNew", data) //GET request to TMTU for routes data
		if err != nil {
			log.Fatal("hi")
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(respRouteNo.Body)
		bodyRouteNo, err := io.ReadAll(respRouteNo.Body) // response body is []byte
		if err != nil {
			fmt.Println("wrong here2")
		}
		var resultRouteNo ResponseRouteNo
		if err := json.Unmarshal(bodyRouteNo, &resultRouteNo); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("wrong here3")
		}
		//fmt.Print(resultRouteNo)
		routes := geojson.NewFeatureCollection()

		for j := 0; j < len(resultRouteNo.Data[0].RouteDetails); j++ {

			sWaypointNo, err := strconv.ParseInt(resultRouteNo.Data[0].RouteDetails[j].Waypoints.WPointNo, 10, 64)
			if err != nil {
				fmt.Println("wrong here15")
			}
			flag := 0
			for k := 0; k < len(stops); k++ {

				if int64(ref[k]) == sWaypointNo {
					flag = 1
				}
			}
			if flag == 0 {
				ref = append(ref, int(sWaypointNo))
			}

			sresultRouteNoLatitude, err := strconv.ParseFloat(resultRouteNo.Data[0].RouteDetails[j].Waypoints.Latitude, 64)
			if err != nil {
				fmt.Println("wrong here16")
			}
			sresultRouteNoLongitude, err := strconv.ParseFloat(resultRouteNo.Data[0].RouteDetails[j].Waypoints.Longitude, 64)
			if err != nil {
				fmt.Println("wrong here1")
			}
			feature := geojson.NewPointFeature([]float64{sresultRouteNoLongitude, sresultRouteNoLatitude})
			feature.SetProperty("name", resultRouteNo.Data[0].RouteDetails[j].Waypoints.WpointName)
			feature.SetProperty("ref", resultRouteNo.Data[0].RouteDetails[j].Waypoints.WPointNo)
			feature.SetProperty("position", j)
			routes.AddFeature(feature)

			for k := 0; k < len(ref); k++ {
				if ref[k] == int(sWaypointNo) {
					stops[int(sWaypointNo)] = resultRouteNo.Data[0].RouteDetails[j].Waypoints.WpointName
					feature1 := geojson.NewPointFeature([]float64{sresultRouteNoLongitude, sresultRouteNoLatitude})
					feature1.SetProperty("name", resultRouteNo.Data[0].RouteDetails[j].Waypoints.WpointName)
					feature1.SetProperty("ref", sWaypointNo)
					feature1.SetProperty("highway", "bus_stop")
					feature1.SetProperty("operator", "Thane Municipal Transport")
					feature1.SetProperty("public_transport", "platform")
					waypoints.AddFeature(feature1)
				}
			}
		}

		rawJSON1, err := routes.MarshalJSON()
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		fn := fmt.Sprintf("output/TMTRoutes%s-%s.json", resultRoutes.Data[i].RouteNo, resultRoutes.Data[i].RouteNum)
		err = os.WriteFile(fn, rawJSON1, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	rawJSON, err := waypoints.MarshalJSON()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//Saves the geojson file for bus stops to the current directory
	//fmt.Printf("%s", string(rawJSON))
	err = os.WriteFile("output/TMTStopsThroughRoutes.json", rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func buslocations() {

	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://tmtu-buslocations-default-rtdb.asia-southeast1.firebasedatabase.app",
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("serviceAccount.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	respBusLocations, err := http.Get("http://tmtitsapi.locationtracker.com/api/getLastTrackingData") //GET request to TMTU for BusLocations data
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(respBusLocations.Body)

	bodyBusLocations, err := io.ReadAll(respBusLocations.Body) // response body is []byte
	if err != nil {
		fmt.Println(err)
	}
	var resultBusLocations ResponseBusLocations
	if err := json.Unmarshal(bodyBusLocations, &resultBusLocations); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	for i := 0; i < len(resultBusLocations.Data); i++ {
		fmt.Printf("On bus no:%s \n", resultBusLocations.Data[i].GetVehicle.Vehid)
		//busID, err := strconv.ParseInt(resultBusLocations.Data[i].GetVehicle.Vehid, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		ref := client.NewRef("tmtu/bus_locations")

		usersRef := ref.Child(resultBusLocations.Data[i].GetVehicle.Vehid)
		if _, errRTDB := usersRef.Push(ctx, &BusLocation{
			Vehid:         resultBusLocations.Data[i].GetVehicle.Vehid,
			VehNo:         resultBusLocations.Data[i].GetVehicle.VehNo,
			RouteNo:       resultBusLocations.Data[i].RouteNo,
			Longitude:     resultBusLocations.Data[i].Longitude,
			Latitude:      resultBusLocations.Data[i].Latitude,
			From:          resultBusLocations.Data[i].DirectionFrom,
			To:            resultBusLocations.Data[i].DirectionTo,
			Waypoint:      resultBusLocations.Data[i].WPointNo,
			LastWaypoint:  resultBusLocations.Data[i].Lastwaypointid,
			LastTrackTime: resultBusLocations.Data[i].LastTrackdt,
			DispatchTime:  resultBusLocations.Data[i].DispatchDateTime,
			CheckedAt:     time.Now().String(),
		}); errRTDB != nil {
			log.Fatalln("Error pushing child node:", errRTDB)
		}

	}

	respLimitRemaining := respBusLocations.Header.Get("X-RateLimit-Remaining")
	respLimitRemainingint, err := strconv.ParseInt(respLimitRemaining, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	for {
		if 2 < respLimitRemainingint {
			break
		} else {
			time.Sleep(10 * time.Second)
			fmt.Print("Waiting for 5 seconds for the API limit to reset\n")
		}
	}
	fmt.Printf("Entered Data for Bus Locations at %s \n", time.Now())
	fmt.Printf("API Limit Remaining: %s \n", respLimitRemaining)
	fmt.Println("Waiting for 5secs...")
	time.Sleep(5 * time.Second)
}
