package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/windingtheropes/fs.alacriware/auth"
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
	return strings.Replace(path, os.Getenv("PUBDIR"), "", -1)
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
	var list string = safe_path(path)
	if len(files) > 250 {
		return fmt.Sprintf("Too many files (%v)", len(files)), nil
	}
	for i := 0; i < len(files); i++ {
		list = list + "\n" + safe_path(filepath.Join(path, files[i].Name()))
	}
	return list, nil
}
func main() {
	// os.Setenv("PUBDIR", "./public")
	public_path := os.Getenv("PUBDIR")
	if public_path == "" {
		fmt.Println("No value found for PUBDIR.")
		os.Exit(1)
	}
	// initialize router
	r := gin.Default()
	r.Use(auth.Auth())
	r.SetTrustedProxies(nil)

	// all paths are registered and checked as routes
	r.NoRoute(func(c *gin.Context) {
		full_path := filepath.Join(public_path, c.Request.URL.String())
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
				c.String(200, list)
				return
			}
		}
		c.Status(404)
	})
	r.Run(":3030")
}
