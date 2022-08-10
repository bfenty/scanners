package main

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
"fmt"
"regexp"
"os"
// "net/http"
// "strings"
)

//Authenticate user from DB
func userauth(userid string) (message Message){
    // Get a database handle.
    var err error
    var user string
    fmt.Println("user:",os.Getenv("USER"))
    fmt.Println("pass:",os.Getenv("PASS"))
    fmt.Println("server:",os.Getenv("SERVER"))
    fmt.Println("port:",os.Getenv("PORT"))

    reg, err := regexp.Compile("[^0-9]+")
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      return message
    }
    // var skus []Sku
    connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders"
    db, err := sql.Open("mysql",
		connectstring)
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      return message
    }

    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
      message.Success = false
      message.Message = pingErr.Error()
      return message
    }
    //set Variables
    //Query
    var newquery string = "select username from users where usercode = ?"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,reg.ReplaceAllString(userid, ""))
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      return message
    }
    defer rows.Close()
    //Pull Data
    for rows.Next() {
    	err := rows.Scan(&user)
    	if err != nil {
        message.Success = false
        message.Message = err.Error()
        return message
    	}
    }
    err = rows.Err()
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      return message
    }
	if user == "" {
    message.Success = false
    message.Message = "User not found. Please scan again."
		return message
	}

  message.Success = true
  message.User = user
  message.Message = "Success"
  return message
}

//Authenticate user from DB
func insert(user string, ordernum string, station string, override bool) (message Message){
    //// DEBUG:
    fmt.Println("USER:"+user+" ORDER:"+ordernum)
    // var message Message

    // Get a database handle.
    var err error

    //open the database
    connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders"
    db, err := sql.Open("mysql",
		connectstring)
    if err != nil {
        // log.Fatal(err)
        message.Success = false
        message.Message = err.Error()
        return message
    }

    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
        message.Success = false
        message.Message = err.Error()
        return message
    }
    if override == false {
    //check if it's already Inserted
    var testquery string = "SELECT COUNT(*) from scans where ordernum = ? AND station = ?"
    rows2, err := db.Query(testquery,ordernum,station)
    if err != nil {
    	log.Fatal(err)
      message.Success = false
      message.Message = err.Error()
      return message
    }
    var val uint
    if rows2.Next() {
      rows2.Scan(&val)
    }
    if(val > 0) {
      fmt.Println(val)
      fmt.Println("Already Entered")
      message.Success = false
      message.Message = "This order has already been scanned. Would you like to override?"
      return message
    }
    }
    //set Variables
    //Query
    var newquery string = "INSERT INTO `scans`(`user`,`ordernum`,`station`) VALUES (?,?,?)"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,user,ordernum,station)
    if err != nil {
    	log.Fatal(err)
      message.Success = false
      message.Message = err.Error()
      return message
    }
    err = rows.Err()
    if err != nil {
    	log.Fatal(err)
      message.Success = false
      message.Message = err.Error()
      return message
    }

    message.Success = true
    message.Message = "Success"
    return message
}
