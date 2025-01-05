package based

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB Database

type Group struct {
	ID   int
	Name string
}
type User struct {
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
	Resource_Path   string
	Group_ID        int
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
type Membership struct {
	ID int
	User_ID int
	Group_ID int
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

// find a single user by id
func (d *Database) GetUser(id int) ([]User, error) {
	// A users slice to hold data from returned rows.
	var users []User

	rows, err := d.db.Query("SELECT * FROM usr WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("find user %q: %v", id, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.ID, &usr.Name); err != nil {
			return nil, fmt.Errorf("user %q: %v", id, err)
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get tokens %q: %v", id, err)
	}
	return users, nil
}

// find the membership of a user to groups, by user id, returns a slice of group ids
func (d *Database) GetUserMembership(user_id int) ([]int, error) {
	// A users slice to hold data from returned rows.
	var memberships []Membership

	rows, err := d.db.Query("SELECT * FROM membership WHERE user_id = ?", user_id)
	if err != nil {
		return nil, fmt.Errorf("error find user memberships %q: %v", user_id, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var memb Membership
		if err := rows.Scan(&memb.ID, &memb.User_ID, &memb.Group_ID); err != nil {
			return nil, fmt.Errorf("error scanning membership to slice, uid %q: %v", user_id, err)
		}
		memberships = append(memberships, memb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over membership, uid %q: %v", user_id, err)
	}

	var group_ids []int
	for i := 0; i < len(memberships); i++ {
		group_ids = append(group_ids, memberships[i].Group_ID)
	}
	return group_ids, nil
}

// get all permissions rows related to a group id
func (d *Database) GetPermissions(group_id int) ([]Permissions, error) {
	// A users slice to hold data from returned rows.
	var permissions []Permissions

	rows, err := d.db.Query("SELECT * FROM permissions WHERE group_id = ?", group_id)
	if err != nil {
		return nil, fmt.Errorf("Error querying permissions %q: %v", group_id, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var perm Permissions
		if err := rows.Scan(&perm.ID, &perm.Resource_Path, &perm.Group_ID, &perm.Allowed, &perm.Apply_Recursive); err != nil {
			return nil, fmt.Errorf("error scanning permissions to slice %q: %v", group_id, err)
		}
		permissions = append(permissions, perm)
	}
	// make sure no errors arose
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions %q: %v", group_id, err)
	}

	return permissions, nil
}


func (d *Database) LogRequest(request Request_Log) (int64, error) {
	result, err := d.db.Exec("INSERT INTO request (ip, access_time, resource_path, token, code) VALUES (?, ?, ?, ?, ?)", request.IP, request.Access_Time, request.Resource_Path, request.Token, request.Code)
	if err != nil {
		return 0, fmt.Errorf("logRequest: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("logRequest: %v", err)
	}
	return id, nil
}

func (d *Database) UpdateToken(token Token) (int64, error) {
	result, err := d.db.Exec("UPDATE token SET user_id = ?, expiry = ?, max = ?, used = ? WHERE id = ?", token.User_ID, token.Expiry, token.Max, token.Used, token.ID)
	if err != nil {
		return 0, fmt.Errorf("update token: %v", err)
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update token: %v", err)
	}
	return rowsUpdated, nil
}

// func (d *Database) AddUser(usr Usr) (int64, error) {
// 	result, err := d.db.Exec("INSERT INTO usr (usr_name) VALUES (?)", usr.Name)
// 	if err != nil {
// 		return 0, fmt.Errorf("add user: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("add user: %v", err)
// 	}
// 	return id, nil
// }

// func (d *Database) AddGrp(grp Grp) (int64, error) {
// 	result, err := d.db.Exec("INSERT INTO grp (grp_name) VALUES (?)", grp.Name)
// 	if err != nil {
// 		return 0, fmt.Errorf("add grp: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("add grp: %v", err)
// 	}
// 	return id, nil
// }

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
