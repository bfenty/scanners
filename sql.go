package main

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
// "log"
"fmt"
"regexp"
"os"
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
      fmt.Println("Message: ",message.Message)
      return message
    }
    connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders"
    db, err := sql.Open("mysql",
		connectstring)
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }

    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
      message.Success = false
      message.Message = pingErr.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
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
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
    defer rows.Close()
    //Pull Data
    for rows.Next() {
    	err := rows.Scan(&user)
    	if err != nil {
        message.Success = false
        message.Message = err.Error()
        fmt.Println("Message: ",message.Message)
        db.Close()
        return message
    	}
    }
    err = rows.Err()
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
	if user == "" {
    message.Success = false
    message.Message = "User not found. Please scan again."
    fmt.Println("Message: ",message.Message)
    db.Close()
		return message
	}

  message.Success = true
  message.User = user
  message.Message = "Success"
  fmt.Println("Message: ",message.Message)
  db.Close()
  return message
}

//Authenticate user from DB
func insert(user string, ordernum string, station string, override bool) (message Message){
    //// DEBUG:
    fmt.Println("USER:"+user+" ORDER:"+ordernum)
    reg, err := regexp.Compile("[^0-9]+")
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      return message
    }
    if len(ordernum) != len(reg.ReplaceAllString(ordernum, "")) {
      message.Success = false
      message.Message = "This doesn't appear to be a valid order id, please scan again"
      fmt.Println("Message: ",message.Message)
      return message
    }

    // Get a database handle.
    // var err error

    //open the database
    connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders"
    db, err := sql.Open("mysql",
		connectstring)
    if err != nil {
        message.Success = false
        message.Message = err.Error()
        fmt.Println("Message: ",message.Message)
        db.Close()
        return message
    }

    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
        message.Success = false
        message.Message = err.Error()
        fmt.Println("Message: ",message.Message)
        db.Close()
        return message
    }

    //Check if the order number is picked before shipping
    if (override == false) && (station=="ship") {
    fmt.Println("Checking if the order was picked...")
    var testquery string = "SELECT COUNT(*) from scans where ordernum = ? AND station = ?"
    rows2, err := db.Query(testquery,ordernum,"pick")
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
    var val uint
    if rows2.Next() {
      rows2.Scan(&val)
    }
    if(val == 0) {
      fmt.Println(val)
      fmt.Println("Order being shipped but not yet picked.")
      message.Success = false
      message.Message = "This order has not yet been picked. Would you like to override?"
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
    }

    //Check if the order number is already inserted
    if override == false {
    fmt.Println("Checking if the order already exists...")
    var testquery string = "SELECT COUNT(*) from scans where ordernum = ? AND station = ?"
    rows2, err := db.Query(testquery,ordernum,station)
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
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
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
    }
    //set Variables
    //Query
    var newquery string = "REPLACE INTO `scans`(`user`,`ordernum`,`station`) VALUES (?,?,?)"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,user,ordernum,station)
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }
    err = rows.Err()
    if err != nil {
      message.Success = false
      message.Message = err.Error()
      fmt.Println("Message: ",message.Message)
      db.Close()
      return message
    }

    message.Success = true
    message.Message = "Success"
    fmt.Println("Message: ",message.Message)
    db.Close()
    return message
}
