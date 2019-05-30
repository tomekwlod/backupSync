package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/tomekwlod/utils"
	"github.com/tomekwlod/utils/sftp"
)

// go run . -path="./" -mustcompile="\\.csv$" -location=chainsaw-backup

func getSFTPLocation(key string) (location *Location, err error) {

	jsonFile, err := os.Open("./locations.json")
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

func main() {

	fmt.Printf("\n")

	fpath := flag.String("path", "./", "Local path")
	fmustcompile := flag.String("mustcompile", "", "Only some files? Usage: \\.csv$ ")
	flocation := flag.String("location", "chainsaw", "One of the sftp locations; It must exist in yml file")
	flag.Parse()

	path := *fpath
	mustcompile := *fmustcompile
	location := *flocation

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

	client, err := sftp.NewClient(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Printf("\nSending files\n")
	for _, file := range files {

		_, err := client.Lstat("/files/" + file.Name)
		if err == nil {
			fmt.Printf("------ File %s already exists at location `%s`\n", file.Name, location)
			continue
		}

		size := math.Round((float64(file.Size)/1024)*100) / 100

		fmt.Printf("-----> Sending file to location `%s`\t%s\t[size:%.2fKB, date:%s]\n", location, file.Filepath, size, file.Time)

		bytesSent, err := client.SendFile("./", file.Name, "/files/", file.Name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("++++++ %d bytes sent\n\n", bytesSent)
	}

	fmt.Printf("\n\nAll done\n\n")
}
