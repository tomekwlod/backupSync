package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/tomekwlod/utils"
	ml "github.com/tomekwlod/utils/logger"
	"github.com/tomekwlod/utils/sftp"
)

// local
// go run . -path="./" -mustcompile="\\.csv$" -location=chainsaw-backup
// prod debug
// go run main.go models.go -dryrun -location="chainsaw-backup" -mustcompile="\\.go$"

var l *ml.Logger

func init() {
	// definig the logger & a log file
	logfile := "/var/log/backupreport.log"

	fmt.Println("LOGFILE:" + logfile)
	fmt.Println("To change the log do: export BACKUPLOGPATH=/var/log/backupsync.log")
	fmt.Println()
	fmt.Println("Usage example:")
	fmt.Println(`backupsync -location="chainsaw-backup" -mustcompile="\\.go$" -dryrun`)

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file" + logfile)
	}

	multi := io.MultiWriter(file, os.Stdout)
	l = ml.New(
		os.Getenv("LOGGING_MODE"),
		// "DEBUG",
		log.New(multi, "", log.Ldate|log.Ltime),
	)
}

func main() {

	fmt.Printf("\n")

	fpath := flag.String("path", "./", "Local path")
	fmustcompile := flag.String("mustcompile", "", "Only some files? Usage: \\.csv$ ")
	flocation := flag.String("location", "chainsaw", "One of the sftp locations; It must exist in yml file")
	fdryrun := flag.Bool("dryrun", false, "For testing, -dryrun or -dryrun=true for Positive, -dryrun=false or -dryrun=0 for Negative")
	flag.Parse()

	path := *fpath
	mustcompile := *fmustcompile
	location := *flocation
	dryrun := *fdryrun

	sftpLocation, err := getSFTPLocation(location)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Scanning %s directory\n", path)
	files := utils.FilesFromDirectory(path, mustcompile)

	if len(files) == 0 {
		fmt.Println("No files to upload found")

		return
	}

	fmt.Printf("\nFound %d file(s)\n", len(files))

	fmt.Printf("\nConnecting to sftp location `%s`\n", location)
	config := &sftp.ClientConfig{
		Username: sftpLocation.Auth.Username,
		Password: sftpLocation.Auth.Password,
		Host:     sftpLocation.Host,
		Port:     sftpLocation.Port,
	}

	backupPath := filepath.Join(sftpLocation.Basepath, "backup")

	client, err := sftp.NewClient(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Printf("\nSending files\n")
	for _, file := range files {

		_, err := client.Lstat(filepath.Join(backupPath, file.Name))
		if err == nil {
			fmt.Printf("------ File %s already exists at location `%s`\n", file.Name, location)
			continue
		}

		size := math.Round((float64(file.Size)/1024)*100) / 100

		fmt.Printf("-----> Sending file to location `%s`\t%s\t[size:%.2fKB, date:%s]\n", location, filepath.Join(backupPath, file.Filepath), size, file.Time)

		if dryrun {
			fmt.Printf("------ DRYRUN is on\n\n")

			continue
		}

		bytesSent, err := client.SendFile("./", file.Name, backupPath, file.Name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("++++++ %d bytes sent\n\n", bytesSent)

		l.Printf("[Location: `%s`] %sbytes %d\t%s", location, file.Time.Format("2006/01/02 15:04:05"), file.Size, filepath.Join(backupPath, file.Filepath))
	}

	fmt.Printf("\n\nAll done\n\n")
}

func getSFTPLocation(key string) (location *Location, err error) {

	jsonFile, err := os.Open("locations.json")
	if err != nil {
		return
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	var locations []*Location
	err = json.Unmarshal(byteValue, &locations)
	if err != nil {
		return
	}

	for i := range locations {
		if locations[i].Name == key {
			location = locations[i]

			return
		}
	}

	return location, errors.New("Location couldn't be determined")
}
