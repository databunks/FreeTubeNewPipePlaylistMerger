# FreetubeNewPipePlaylistMerger
A simple Golang script to merge FreeTube playlists into NewPipe

## Why this was created
While there are tools to merge NewPipe playlists into FreeTube, as of the time of the creation of this script there has been no known scripts which merge FreeTube playlists into NewPipe,
so this has been created as an easier method of doing so.

## How to use
1) cd into the same directory as main.go
2) Run the command:
`go run main.go -f "<path to freetube playlists file (.db)>" -n "<Path to newpipe database file"`

## Optional
- It is highly recommended to use this with syncthing (https://syncthing.net/) for ease of use, that way you can export to a shared folder and easily transfer files.

## Important Notes
- Must have Golang installed in order for it to run. (`sudo apt install go` `sudo pacman -S go`)
- This script is not designed to handle playlists with duplicate names so please make sure each playlist has a unique name.
- Thumbnails for playlists will first appear as empty as the script is designed to work when offline but they will eventually load when you click into them, you can reset the thumbnail of the playlist inside of NewPipe too. (feel free to make a pull request to add an online version)
- This works best on Linux. (has not been tested on windows)
