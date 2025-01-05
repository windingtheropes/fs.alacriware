package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/windingtheropes/fs.alacriware/auth"
	"github.com/windingtheropes/fs.alacriware/based"
	"github.com/windingtheropes/fs.alacriware/based/webdb"
	"github.com/windingtheropes/fs.alacriware/logger"
)

// check if path exists
func path_exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// check is path a file
func is_file(path string) bool {
	info, err := os.Stat(path)
	// if any error return false, but this should always be run after checking if exists
	if os.IsNotExist(err) || err != nil {
		fmt.Println("Error with path " + path + ":" + err.Error())
		return false
	}
	return !info.IsDir()
}

func safe_path(path string) string {
	// Windows is dumb
	path = strings.ReplaceAll(path, "\\", "/")
	cleaned_path := strings.Replace(path, os.Getenv("PUBDIR"), "", 1)
	if cleaned_path == "" {
		return "/"
	} else {
		return cleaned_path
	}
}
func get_dir_list(path string) (string, error) {
	info, err := os.Stat(path)
	// if any error return false, but this should always be run after checking if exists
	if os.IsNotExist(err) || err != nil {
		fmt.Println("Error with path " + path + ":" + err.Error())
		return "", err
	}
	if !info.IsDir() {
		fmt.Println(path + " is not a directory.")
		return "", err
	}
	if err != nil {
		fmt.Println("Error while gen dir path " + path + ":" + err.Error())
		return "", err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error while gen dir path " + path + ":" + err.Error())
		return "", err
	}

	var list string = fmt.Sprintf("<h1>%v</h1>",safe_path(path))
	if len(files) > 250 {
		return fmt.Sprintf("Too many files (%v)", len(files)), nil
	}
	for i := 0; i < len(files); i++ {
		rel_path := safe_path(filepath.Join(path, files[i].Name()))
		list = list + "<br>" + fmt.Sprintf("<a href='%v'>%v</>",rel_path,rel_path)
	}
	return list, nil
}
func main() {
	// go admin_server()
	file_server()
}
func file_server() {
	public_path := os.Getenv("PUBDIR")
	if public_path == "" {
		fmt.Println("No value found for PUBDIR.")
		os.Exit(1)
	}
	// initialize router
	r := gin.Default()

	config := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", os.Getenv("DBHOST"), os.Getenv("DBPORT")),
		DBName: os.Getenv("DBNAME"),
	}

	webdb.Init(based.ConectDB(config))
	
	r.Use(auth.Auth())
	r.Use(logger.LogRequest())
	r.SetTrustedProxies(nil)

	// all paths are registered and checked as routes
	r.NoRoute(func(c *gin.Context) {
		full_path := filepath.Join(public_path, c.Request.URL.Path)
		if path_exists(full_path) {
			if is_file(full_path) {
				// Is file
				c.File(full_path)
				return
			} else {
				// Is directory
				list, err := get_dir_list(full_path)
				if err != nil {
					c.Status(500)
					return
				}
				c.Data(200, "text/html; charset=utf-8", []byte(list))
				return
			}
		}
		c.Status(404)
	})
	r.Run(":3030")
}
// func admin_server() {
// 	// initialize router
// 	r := gin.Default()

// 	r.GET("/hello", func(c *gin.Context) {
// 		c.Status(200)
// 	})
// 	r.Run(":4040")
// }
