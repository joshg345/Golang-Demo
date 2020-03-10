package main

import "github.com/go-resty/resty"
import "github.com/buger/jsonparser"
import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

type Title struct {
	Title string `gorm:"unique;not null"`
}

type Stat struct {
	DBN          string
	SchoolName   string
	NoTests      string
	ReadingMeans string
	MathsMeans   string
	WriteMeans   string
}

func main() {
	// Create a Resty Client
	client := resty.New()
	resp, _ := client.R().
		EnableTrace().
		Get("https://data.cityofnewyork.us/api/views/zt9s-n5aj/rows.json")

	// output JSON data to slice variable
	var data []byte = resp.Body()

	//*** Commented out this section as no longer needed but left in for reference

	// var fieldName []string

	// jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	//Get data using jsonparser first returned value
	// 	var id, _ = jsonparser.GetInt(value, "id")
	// 	var name, _ = jsonparser.GetString(value, "name")
	// 	if id != -1 {
	// 		fieldName = append(fieldName, name)
	// 	}
	// }, "meta", "view", "columns")

	//declare slice variables
	var DBNs []string
	var sclNames []string
	var numTests []string
	var readingMeans []string
	var mathsMeans []string
	var writeMeans []string

	// move each piece of data into a slice array
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//Get data using jsonparser first returned value
		var DBN, _ = jsonparser.GetString(value, "[8]")
		var sclName, _ = jsonparser.GetString(value, "[9]")
		var numTest, _ = jsonparser.GetString(value, "[10]")
		var readingMean, _ = jsonparser.GetString(value, "[11]")
		var mathsMean, _ = jsonparser.GetString(value, "[12]")
		var writeMean, _ = jsonparser.GetString(value, "[13]")

		DBNs = append(DBNs, DBN)
		sclNames = append(sclNames, sclName)
		numTests = append(numTests, numTest)
		readingMeans = append(readingMeans, readingMean)
		mathsMeans = append(mathsMeans, mathsMean)
		writeMeans = append(writeMeans, writeMean)

	}, "data")

	// create database connection
	db, err := gorm.Open("sqlite3", "schoolstats.db")
	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	// delete old tables
	// db.Delete(&Title{})
	db.Delete(&Stat{})

	// db.AutoMigrate(&Title{})
	db.AutoMigrate(&Stat{})

	// for i := 0; i < len(fieldName); i++ {
	// 	var titles = Title{Title: fieldName[i]}
	// 	db.Create(&titles)
	// }

	// for each entry in array, create an entry in Stats table
	for i := 0; i < len(DBNs); i++ {
		var stats = Stat{DBN: DBNs[i], SchoolName: sclNames[i], NoTests: numTests[i], ReadingMeans: readingMeans[i], MathsMeans: mathsMeans[i], WriteMeans: writeMeans[i]}
		db.Create(&stats)
	}

}
