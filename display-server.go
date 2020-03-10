package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/labstack/echo"
  "net/http"
)

// User
type Data struct {
  DBN          string `json:"dbn"`
  SchoolName   string `json:"school_name"`
  NoTests      string `json:"no_tests"`
  ReadingMeans string `json:"reading_means"`
  MathsMeans   string `json:"maths_means"`
  WriteMeans   string `json:"write_means"`
}

func main() {
  // Initialise echo
  e := echo.New()
  // Set up home routing
  e.GET("/", getData)
  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}

func getData(c echo.Context) error {
  // open conncetion to database
  db, err := gorm.Open("sqlite3", "schoolstats.db")
  if err != nil {
    panic("failed to connect database")
  }

  defer db.Close()
  // variable to store table data
  var datas []Data
  // get all data in table and output to slice
  db.Table("stats").Select("*").Scan(&datas)
  // return JSON data to server
  return c.JSON(http.StatusOK, &datas)
}
