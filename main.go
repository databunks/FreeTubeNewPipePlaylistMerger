package main

import (
	"fmt"
	"flag"
	"encoding/json"
	"os"
	"bufio"
	"database/sql"
	_ "modernc.org/sqlite"
	"playlistMerger/database"
	"playlistMerger/zip"
)

type Video struct {
	VideoID 	string 	`json: "videoId"`
	Title   	string	`json: "title"`
	Author	        string	`json: "author"`
	AuthorId	string	`json: "authorId"`
	LengthSeconds	int	`json: "lengthSeconds"`
	TimeAdded	int	`json: "timeAdded"`
	PlaylistItemId	string	`json: "playlistItemId"`
	Type		string	`json: "type"`
}

type Playlist struct {
	PlaylistName 	string 	  `json: "playlistName"`
	Description  	string 	  `json: "description"`
	Videos 	     	[]Video   `json: "videos"`
	ID	     	string	  `json: "_id"`
	CreatedAt    	int	  `json: "createdAt"`
	LastUpdatedAt	int	  `json: "lastUpdatedAt"`		
}

func ReadFreeTubePlaylists(path string) []Playlist {
	var playlists []Playlist

	// Read playlists from freetube
	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var playlist Playlist
		json.Unmarshal([]byte(scanner.Text()), &playlist)
		playlists = append(playlists, playlist)
	}

	return playlists
}

func MergePlaylists(playlists []Playlist, newPipeDBPath string){
	// Merge with newpipe .db file
	db, err := sql.Open("sqlite", newPipeDBPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	databaseRepository := &database.DatabaseRepository{Db: db}
	
	for _, playlist := range playlists {
		// Add playlist in if there is none with a matching name
		if (!databaseRepository.CheckPlaylistExists(playlist.PlaylistName)){
			databaseRepository.InsertPlaylist(database.Playlist{UID: databaseRepository.GetMaxPlaylistUID() + 1, Name: playlist.PlaylistName, IsThumbnailPermanent: 0, ThumbnailStreamID: 1, DisplayIndex: databaseRepository.GetMaxPlaylistDisplayIndex() + 1 })
			fmt.Printf("Created Playlist %s\n", playlist.PlaylistName)
		}

		for _, video := range playlist.Videos {
			streamURL := "https://www.youtube.com/watch?v=" + video.VideoID
			stream := database.Stream{UID: databaseRepository.GetMaxStreamUID() + 1, ServiceID: 0, URL: streamURL, Title: video.Title, StreamType: "VIDEO_STREAM", Duration: video.LengthSeconds, Uploader: "", UploaderURL: "", ThumbnailURL: "", ViewCount: 0, TextualUploadDate: "", UploadDate: 0, IsUploadDateApprox: 1}

			if (!databaseRepository.CheckStreamExists(streamURL)){
				databaseRepository.InsertStream(stream)
			} else {
				stream = databaseRepository.GetStream(streamURL)
			}

			databaseRepository.InsertStreamIntoPlaylist(stream, playlist.PlaylistName)
		}
	}

}

func main() {
	freeTubePlaylistsFilePath := flag.String("f", "", "Path to freetube playlists file (.db)")
	newPipeDatabaseFolderFilePath := flag.String("n", "", "Path to newpipe database file (.zip)")

	flag.Parse()

	if (*freeTubePlaylistsFilePath == ""){
		panic("You must provide a file path to your freetube playlists folder with -f \"pathtofile\"")
	}

	if (*newPipeDatabaseFolderFilePath == ""){
		panic("You must provide a file path to your newpipe database folder with -n \"pathtofile\"")
	}


	fmt.Println(*freeTubePlaylistsFilePath)
	fmt.Println(*newPipeDatabaseFolderFilePath)

	// Read playlists from freetube
	fmt.Println("Reading From Freetube Playlists File....")
	playlists := ReadFreeTubePlaylists(*freeTubePlaylistsFilePath)	
	
	// Unzip database folder
	fmt.Println("Unzipping New Pipe Folder....")
	zip.DecompressFolder(*newPipeDatabaseFolderFilePath)

	// Merge playlists
	fmt.Println("Merging Into Newpipe Databse....")
	MergePlaylists(playlists, "./output/newpipe.db")

	// Zip database folder
	fmt.Println("Zipping merged database....")
	zip.CompressFiles("./output")

	// Cleanup leftover files
	fmt.Println("Cleaning up....")
	zip.RemoveFolder("./output")

	fmt.Println("Finished Merging!")
}


