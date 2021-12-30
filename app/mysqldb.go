package main

import (
	"database/sql"
	"fmt"
	"os"
)

var (
	// database information
	dbhostsip  = os.Getenv("DB_HOST")     // IP address of database server
	dbusername = os.Getenv("DB_USERNAME") // username of the database user
	dbpassword = os.Getenv("DB_PASSWORD") // password of the database user
	dbname     = os.Getenv("DB_NAME")     // name of the database
	dbcharset  = os.Getenv("DB_CHARSET")  // database character set

	// list of relations in database needed for querying
	sampleTable = "sample" // table name for airplane
)

// function to generate database connection string
func getConnectionString() string {
	// final output sample: root:password@tcp(127.0.0.1:3306)/database_name?charset=utf-8
	return dbusername + ":" + dbpassword + "@tcp(" + dbhostsip + ")/" + dbname + "?charset=" + dbcharset
}

// function to return JSON array from a MySQL database
func changeDBDataToJSON(sqlString string) ([]map[string]interface{}, error) {
	sqlConnString := getConnectionString()
	db, err := sql.Open("mysql", sqlConnString) //
	if err != nil {
		return nil, err // return error if present
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db) // close database connection if error occurs
	rows, err := db.Query(sqlString) // typically, returns rows querying the database
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns() // returns the column name and returns errors if rows are closed
	if err != nil {
		return nil, err
	}
	count := len(columns) // number of columns from where rows are rendered.
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtr := make([]interface{}, count)
	// JSON format builder
	for rows.Next() { // for each row
		for i := 0; i < count; i++ {
			valuePtr[i] = &values[i] // find value pointers
		}
		err := rows.Scan(valuePtr...)
		if err != nil {
			return nil, err
		}
		entry := make(map[string]interface{})
		for i, col := range columns { // for each data in row build a json data (with key and value4)
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v // assign data as entry["key"] = value
		}
		tableData = append(tableData, entry) // append the data object map[string]interface{} to array []map[string]interface{}
	}
	// for marshalling (i.e. serializing)
	// jsonData, err := json.Marshal(tableData)
	// if err != nil {
	// 	return jsonData, err
	// }

	return tableData, nil
}

// function to get data from MySQL database ==> gets you ready response data
// sample buildout query: SELECT * FROM table_name WHERE id = id_requested;
func getDataDBbyIndex(table string, index string, id string) (int, map[string]interface{}) {
	status := 200
	response := make(map[string]interface{})
	sqlQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s';", table, index, id)
	dbData, err := changeDBDataToJSON(sqlQuery)
	if err != nil {
		status = 500
		response = statMsgData[3] // Error found querying with the database
	} else if len(dbData) == 0 {
		status = 500
		response = statMsgData[4] // Data not found in the database
	} else {
		response = statMsgData[2] // Everything's fine
		response["data"] = dbData
	}
	return status, response
}
