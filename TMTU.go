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
	"path/filepath"
	"strconv"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		City             interface{} `json:"City,omitempty"`
		Speed            string      `json:"Speed"`
		ImagePath        interface{} `json:"ImagePath,omitempty"`
		AC               string      `json:"AC"`
		Ignition         string      `json:"Ignition"`
		AUX1             string      `json:"AUX1"`
		DI4              string      `json:"DI4"`
		Fuel             string      `json:"Fuel"`
		Temparature      string      `json:"Temparature"`
		WPointNo         int         `json:"WPointNo,omitempty"`
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
		Routeflag1       int         `json:"routeflag1,omitempty"`
		RouteNo          string      `json:"RouteNo"`
		WaybillNo        string      `json:"WaybillNo"`
		Lastwaypointid   int         `json:"lastwaypointid,omitempty"`
		Token            string      `json:"token"`
		Avgspeed         string      `json:"avgspeed"`
		LatLong          string      `json:"LatLong"`
		GetVehicle       struct {
			Vehid string `json:"vehid"`
			VehNo string `json:"VehNo"`
		} `json:"get_vehicle"`
	} `json:"data"`
}

type Data struct {
	IdxTrackidPk     int                `bson:"idx_Trackid_pk"`
	VehID            int                `bson:"VehId"`
	VehNo            string             `bson:"VehNo"`
	CmpID            int                `bson:"CmpId"`
	LastTrackdt      primitive.DateTime `bson:"LastTrackdt"`
	NCSent           string             `bson:"NCSent,omitempty"`
	CSent            string             `bson:"CSent,omitempty"`
	PrevTrackDt      primitive.DateTime `bson:"PrevTrackDt"`
	LastNCSentDate   primitive.DateTime `bson:"LastNCSentDate"`
	City             interface{}        `bson:"City,omitempty"`
	Speed            float64            `bson:"Speed"`
	ImagePath        interface{}        `bson:"ImagePath,omitempty"`
	AC               bool               `bson:"AC"`
	Ignition         bool               `bson:"Ignition"`
	AUX1             bool               `bson:"AUX1"`
	DI4              bool               `bson:"DI4"`
	Fuel             float64            `bson:"Fuel,omitempty"`
	Temparature      string             `bson:"Temperature,omitempty"`
	WPointNo         int                `bson:"WPointNo,omitempty"`
	Odometer         float64            `bson:"Odometer"`
	Distance         float64            `bson:"Distance"`
	ETATime          float64            `bson:"ETATime"`
	ETARoute         string             `bson:"ETARoute,omitempty"`
	ETAOldTime       float64            `bson:"ETAOldTime"`
	Routeflag        bool               `bson:"routeflag"`
	ETARouteName     string             `bson:"ETARouteName"`
	DirectionFrom    string             `bson:"DirectionFrom"`
	DirectionTo      string             `bson:"DirectionTo"`
	DispatchDateTime primitive.DateTime `bson:"DispatchDateTime"`
	ETATime1         float64            `bson:"ETATime1"`
	ETAOldTime1      float64            `bson:"ETAOldTime1"`
	Routeflag1       int                `bson:"routeflag1,omitempty"`
	RouteNo          int                `bson:"RouteNo"`
	WaybillNo        int                `bson:"WaybillNo"`
	Lastwaypointid   int                `bson:"lastwaypointid,omitempty"`
	Token            int                `bson:"token"`
	Avgspeed         float64            `bson:"avgspeed"`
	Location         Location
}

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

func main() {
	err := os.MkdirAll("output", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	waypoints()

	routes()
	start := time.Now()
	fmt.Printf("Started At:%s\n", start.String())

	for i := 0; i < 18000; i++ {
		buslocations(i)
		i++
	}
	fmt.Printf("Finished Recording In:%s", time.Since(start))

	fmt.Println("Saving Recording And Parsing to MongoDB:")
	mongoBusLocationSaver(1, 18056)
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
	err := os.MkdirAll("output/rawroutes/routes", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

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

		fn := fmt.Sprintf("output/rawroutes/routes/TMTRoutes%s-%s.json", resultRoutes.Data[i].RouteNo, resultRoutes.Data[i].RouteNum)
		err = os.WriteFile(fn, bodyRouteNo, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func buslocations(i int) {

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

	err1 := os.MkdirAll("output/buslocations", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err1)
	}

	fmt.Printf("Iterating %d times\n", i)
	saveSrc := fmt.Sprintf("output/buslocations/%d.json", i)
	err = os.WriteFile(saveSrc, bodyBusLocations, 0644)
	if err != nil {
		log.Fatal(err)
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
	fmt.Printf("Saved Data for Bus Locations at %s \n", time.Now())
	fmt.Printf("API Limit Remaining: %s \n", respLimitRemaining)
	fmt.Println("Waiting for 5secs...")
	time.Sleep(5 * time.Second)
}

func mongoBusLocationSaver(x int, y int) {
	start := time.Now()
	uri := "mongodb://localhost:27017" //monogodb Connection String
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	for i := x; i <= y; i++ {

		fileNo := strconv.Itoa(i)
		fileName := fileNo + ".json"
		filePath := filepath.Join("output", "buslocations", fileName)
		fileData, err := os.ReadFile(filePath)
		fmt.Printf("Reading file %s \n", fileNo)
		var busLocations ResponseBusLocations
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(fileData, &busLocations)
		for j := 0; j < len(busLocations.Data); j++ {
			lastTrackdtTime, _ := time.Parse("2006-01-02 15:04:05", busLocations.Data[j].LastTrackdt)
			lastTrackdtBson := primitive.NewDateTimeFromTime(lastTrackdtTime)
			coll := client.Database("TMTU").Collection(busLocations.Data[j].VehID)
			var result bson.M
			err = coll.FindOne(context.TODO(), bson.D{{Key: "LastTrackdt", Value: lastTrackdtBson}}).Decode(&result)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					prevTrackdtTime, _ := time.Parse("2006-01-02 15:04:05", busLocations.Data[j].PrevTrackDt)
					prevTrackdtBson := primitive.NewDateTimeFromTime(prevTrackdtTime)
					LastNCSentDateTime, _ := time.Parse("2006-01-02 15:04:05", busLocations.Data[j].LastNCSentDate)
					LastNCSentDateBson := primitive.NewDateTimeFromTime(LastNCSentDateTime)
					DispatchDateTime, _ := time.Parse("2006-01-02 15:04:05", busLocations.Data[j].DispatchDateTime)
					DispatchDateBson := primitive.NewDateTimeFromTime(DispatchDateTime)

					iVehID, _ := strconv.Atoi(busLocations.Data[j].VehID)
					iCmpID, _ := strconv.Atoi(busLocations.Data[j].CmpID)
					iToken, _ := strconv.Atoi(busLocations.Data[j].Token)
					iRouteNo, _ := strconv.Atoi(busLocations.Data[j].RouteNo)
					iWaybillNo, _ := strconv.Atoi(busLocations.Data[j].WaybillNo)

					fbusLocationsLatitude, err := strconv.ParseFloat(busLocations.Data[j].Latitude, 64)
					if err != nil {
						fmt.Println(err)
					}
					fbusLocationsLongitude, err := strconv.ParseFloat(busLocations.Data[j].Longitude, 64)
					if err != nil {
						fmt.Println(err)
					}
					fFuel, err := strconv.ParseFloat(busLocations.Data[j].Fuel, 64)
					if err != nil {
						fmt.Println(err)
					}
					fOdometer, err := strconv.ParseFloat(busLocations.Data[j].Odometer, 64)
					if err != nil {
						fmt.Println(err)
					}
					fDistance, err := strconv.ParseFloat(busLocations.Data[j].Distance, 64)
					if err != nil {
						fmt.Println(err)
					}
					fETATime, err := strconv.ParseFloat(busLocations.Data[j].ETATime, 64)
					if err != nil {
						fmt.Println(err)
					}
					fETAOldTime, err := strconv.ParseFloat(busLocations.Data[j].ETAOldTime, 64)
					if err != nil {
						fmt.Println(err)
					}
					fETATime1, err := strconv.ParseFloat(busLocations.Data[j].ETATime1, 64)
					if err != nil {
						fmt.Println(err)
					}
					fETAOldTime1, err := strconv.ParseFloat(busLocations.Data[j].ETAOldTime1, 64)
					if err != nil {
						fmt.Println(err)
					}
					fAvgspeed, err := strconv.ParseFloat(busLocations.Data[j].Avgspeed, 64)
					if err != nil {
						fmt.Println(err)
					}
					fSpeed, err := strconv.ParseFloat(busLocations.Data[j].Speed, 64)
					if err != nil {
						fmt.Println(err)
					}
					var bIgnition bool
					var bAC bool
					var bDI4 bool
					var bAUX1 bool
					if busLocations.Data[j].Ignition == "ON" {
						bIgnition, err = strconv.ParseBool("true")
					}

					if busLocations.Data[j].AUX1 == "ON" {
						bAUX1, err = strconv.ParseBool("true")
					}
					if busLocations.Data[j].DI4 == "ON" {
						bDI4, err = strconv.ParseBool("true")
					}
					if busLocations.Data[j].AC == "ON" {
						bAC, err = strconv.ParseBool("true")
					}
					if err != nil {
						fmt.Println(err)
					}
					bRouteflag, _ := strconv.ParseBool(busLocations.Data[j].Routeflag)

					if busLocations.Data[j].CSent == "0" {
						busLocations.Data[j].CSent = ""
					}
					if busLocations.Data[j].NCSent == "0" {
						busLocations.Data[j].NCSent = ""
					}
					if busLocations.Data[j].Temparature == "0" {
						busLocations.Data[j].Temparature = ""
					}

					location := Location{
						Type:        "Point",
						Coordinates: []float64{fbusLocationsLongitude, fbusLocationsLatitude},
					}
					bus := Data{
						IdxTrackidPk:     busLocations.Data[j].IdxTrackidPk,
						VehID:            iVehID,
						CmpID:            iCmpID,
						LastTrackdt:      lastTrackdtBson,
						NCSent:           busLocations.Data[j].NCSent,
						CSent:            busLocations.Data[j].CSent,
						PrevTrackDt:      prevTrackdtBson,
						LastNCSentDate:   LastNCSentDateBson,
						City:             busLocations.Data[j].City,
						Speed:            fSpeed,
						ImagePath:        busLocations.Data[j].ImagePath,
						AC:               bAC,
						Ignition:         bIgnition,
						AUX1:             bAUX1,
						DI4:              bDI4,
						Fuel:             fFuel,
						Temparature:      busLocations.Data[j].Temparature,
						WPointNo:         busLocations.Data[j].WPointNo,
						Odometer:         fOdometer,
						Distance:         fDistance,
						ETATime:          fETATime,
						ETARoute:         busLocations.Data[j].ETARoute,
						ETAOldTime:       fETAOldTime,
						Routeflag:        bRouteflag,
						ETARouteName:     busLocations.Data[j].ETARouteName,
						DirectionFrom:    busLocations.Data[j].DirectionFrom,
						DirectionTo:      busLocations.Data[j].DirectionTo,
						DispatchDateTime: DispatchDateBson,
						ETATime1:         fETATime1,
						ETAOldTime1:      fETAOldTime1,
						Routeflag1:       busLocations.Data[j].Routeflag1,
						RouteNo:          iRouteNo,
						WaybillNo:        iWaybillNo,
						Lastwaypointid:   busLocations.Data[j].Lastwaypointid,
						Token:            iToken,
						Avgspeed:         fAvgspeed,
						VehNo:            busLocations.Data[j].GetVehicle.VehNo,
						Location:         location,
					}

					coll.InsertOne(context.TODO(), bus)
				}
				fmt.Print("+")
			}
		}
	}
	fmt.Printf("\nTime taken for MongoDB saver :%s", time.Since(start))
}
