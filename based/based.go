package based

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB Database

type Grp struct {
	ID   int
	Name string
}
type Usr struct {
	ID   int
	Name string
}
type Token struct {
	ID      string
	User_ID int
	Expiry  int64 //milli
	Max     int16
	Used    int16
}
type Permissions struct {
	ID              int64
	User_ID         int
	Group_ID        int
	Resource_Path   string
	Allowed         bool
	Apply_Recursive bool
}
type Request_Log struct {
	IP            string
	Access_Time   int64
	Resource_Path string
	Token         string
	Code          int
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

// find a token
func (d *Database) FindToken(t string) ([]Token, error) {
	// A tokens slice to hold data from returned rows.
	var tokens []Token

	rows, err := d.db.Query("SELECT * FROM token WHERE id = ?", t)
	if err != nil {
		return nil, fmt.Errorf("find token %q: %v", t, err)
	}
	defer rows.Close()
	
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var tok Token
		if err := rows.Scan(&tok.ID, &tok.User_ID, &tok.Expiry, &tok.Max, &tok.Used); err != nil {
			return nil, fmt.Errorf("token %q: %v", t, err)
		}
		tokens = append(tokens, tok)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get tokens %q: %v", t, err)
	}
	return tokens, nil
}

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

func (d *Database) AddUser(usr Usr) (int64, error) {
	result, err := d.db.Exec("INSERT INTO usr (usr_name) VALUES (?)", usr.Name)
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}
	return id, nil
}

func (d *Database) AddGrp(grp Grp) (int64, error) {
	result, err := d.db.Exec("INSERT INTO grp (grp_name) VALUES (?)", grp.Name)
	if err != nil {
		return 0, fmt.Errorf("add grp: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add grp: %v", err)
	}
	return id, nil
}

// func (d *Database) AddToken(tok Token) (int64, error) {
// 	result, err := d.db.Exec("INSERT INTO grp (id, user_id, expiry, max, used) VALUES (?)", grp.Name)
// 	if err != nil {
// 		return 0, fmt.Errorf("add grp: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("add grp: %v", err)
// 	}
// 	return id, nil
// }
