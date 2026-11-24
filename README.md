# FreetubeNewPipePlaylistMerger
A simple golang script to merge freetube playlists into newpipe

## Why this was created
While there are tools to merge newpipe playlists into freetube, as of the time of the creation of this script there has been no known scripts which merge freetube playlists into newpipe,
so this has been created as an easier method of doing so.

## How to use
1) cd into the same directory as main.go
2) Run the command:
`go run main.go -f "<path to freetube playlists file (.db)>" -n "<Path to newpipe database file"`

## Important Notes
- Must have golang installed in order for it to run (`sudo apt install go` `sudo pacman -S go`)
- This script is not designed to handle playlists with duplicate names so please make sure each playlist has a unique name.
- Thumbnails for playlists will first appear as empty as the script is designed to work when offline but they will eventually load when you click into them, you can reset the thumbnail of the playlist inside of newpipe too, (feel free to make a pull request to add an online version)
- This works best on linux (has not been tested on windows)
