package db

import (
	"database/sql"
	"fmt"
	"log"

	types "github.com/Oluwatunmise-olat/Horus-Clone/types"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Password 		 string
	UserName     string
	Port         int
	Host         string
	DatabaseName string
	DiscordGuildId string
	DiscordAppId string
	DiscordToken string
}

type Store struct{
	DBInstance *sql.DB
}


func Connect(dsn string) (*Store, error) {
	/*
	 Connect to mysql database
	*/
	dbConn, err :=  sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("Could not connect to Database. ¦")
	}

	return &Store{ DBInstance: dbConn }, nil
}

func (s *Store) AutoMigrateTable() error {
	/*
	  Create request_logger table
	*/
	query := fmt.Sprintf(
		`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				request_body TEXT NUll,
				request_path VARCHAR(1000),
				request_headers TEXT,
				request_method VARCHAR(10),
				request_host VARCHAR(1000),
				response_body TEXT NULL,
				response_status INT,
				req_res_time DOUBLE
			)
		`, "request_logger",
	)
	_, err := s.DBInstance.Exec(query)

	return err
}

func (s *Store) PerformInsert(data *types.RequestLogger) error {
	_, err := s.DBInstance.Query(
		`
			INSERT INTO request_logger 
				(
					request_body, 
					request_path,
					request_headers, 
					request_method,
					request_host, 
					response_body,
					response_status, 
					req_res_time
				)
				VALUES 
				(?, ?, ?, ?, ?, ?, ?, ?)
		`, 
		data.RequestBody, data.RequestPath,data.RequestHeaders, data.RequestMethod, data.RequestHost, data.ResponseBody, data.ResponseStatus, data.ReqResTime)

	if err != nil {
		return err
	}
	
	return nil
}

func (s *Store) GetLogs(limit *int, method string) []*types.RequestLogger {
	var Limit = 20
	if (limit != nil) {
		Limit = *limit
	}
	records, err := s.DBInstance.Query("SELECT request_body, request_path, request_headers, request_method, request_host, response_body, response_status, req_res_time FROM request_logger WHERE request_method LIKE ? ORDER BY request_logger.id DESC LIMIT ?", method, Limit)
	
	if err != nil {
		log.Fatalf("Could not fetch requests: Err >>>>> &v", err)
	}

	defer records.Close()

	data := []*types.RequestLogger{}

	for records.Next(){
		row := new(types.RequestLogger)
		if err := records.Scan(&row.RequestBody, &row.RequestPath, &row.RequestHeaders, &row.RequestMethod, &row.RequestHost, &row.ResponseBody, &row.ResponseStatus, &row.ReqResTime); err != nil {
			log.Fatalf("Could not scan database rows: Err >>>>> &v", err)
		}

		data = append(data, row)
	}

	return data
}