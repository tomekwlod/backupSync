# Backup tool

Yet another tool to send the files to SFTP area

### Building
First make sure you have your locations settings set in `locations.json` file. 

__Download all needed dependencies__
`go get ./...`

__Build  the code__
`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backupsync .`


### Deploy
Currently just run the `make deploy` command on your prod machine. It presumes you have your local OS scripts live at `/opt/scripts/`. It moves the ./backupscript there under ./backupsync name and adds +x so it can be executed. 
It is not perfect because there is a hardcoded path in the ./backupscript file (_please change it later_). 
Now you can run the command from anywhere


### Usage
Usage is rather simple:

__Locally you can run__
`go run . -path="./" -mustcompile="\\.csv$" -location=chainsaw-backup`

__On prod__
`./backupsync -dryrun -path="./" -location="chainsaw-backup" -mustcompile="\\.gz$"`

- `path` sets the local path the files will be searched for
- `mustcompile` regexp to include only the files you're interested in; leave empty for all the files
- `location` remote location to one of the destinations from the locations.json file
- `dryrun` to test the files, nothing here will be sent to sftp  