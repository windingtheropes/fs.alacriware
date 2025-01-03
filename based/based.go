package based

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB Database

type grp struct {
	ID   int64
	name string
}
type usr struct {
	ID   int64
	name string
}
type token struct {
	ID      int64
	user_ID int64
	expiry  int64
	max     int16
	used    int16
}
type Permissions struct {
	ID              int64
	// UID 			int NEED A WAY TO IDENTIFY BOTH GROUPS AND USERS 
	Resource_Path   string
	Allowed         bool
	Apply_Recursive bool
}
type Request_Log struct {
	IP            string
	Access_Time   int64
	Resource_Path string
	Token         string
	Code		  int
}

func InitDB(config mysql.Config) {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Database Connected.")

	DB = Database{
		db,
	}
}

type Database struct {
	db *sql.DB
}

func (d *Database) Hi() {
	fmt.Println("i'm based!")
}

// // albumsByArtist queries for albums that have the specified artist name.
// func (d *Database) albumsByArtist(name string) ([]Album, error) {
//     // An albums slice to hold data from returned rows.
//     var albums []Album

//     rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
//     if err != nil {
//         return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
//     }
//     defer rows.Close()
//     // Loop through rows, using Scan to assign column data to struct fields.
//     for rows.Next() {
//         var alb Album
//         if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
//             return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
//         }
//         albums = append(albums, alb)
//     }
//     if err := rows.Err(); err != nil {
//         return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
//     }
//     return albums, nil
// }

func (d *Database) LogRequest(request Request_Log) (int64, error) {
	result, err := d.db.Exec("INSERT INTO requests (ip, access_time, resource_path, token, code) VALUES (?, ?, ?, ?, ?)", request.IP, request.Access_Time, request.Resource_Path, request.Token, request.Code)
	if err != nil {
		return 0, fmt.Errorf("logRequest: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("logRequest: %v", err)
	}
	return id, nil
}
