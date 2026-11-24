package database

import (
       "database/sql"
       "strconv"
       "fmt"
)


type DatabaseRepository struct {
	Db *sql.DB
}

type Playlist struct {
	UID			int
	Name			string
	IsThumbnailPermanent	int
	ThumbnailStreamID	int
	DisplayIndex		int
}

type Stream struct {
	UID			int
	ServiceID		int
	URL			string
	Title			string
	StreamType		string
	Duration		int
	Uploader		string
	UploaderURL		string
	ThumbnailURL		string
	ViewCount		int
	TextualUploadDate	string
	UploadDate		int
	IsUploadDateApprox	int
}

type PlaylistItem struct {
	PlaylistID 		int
	StreamID 		int
	JoinIndex		int
	UID			int
	serviceID		int
	URL			string
	Title			string
	StreamType		string
	Duration		int
	Uploader		string
	UploaderURL		string
	ThumbnailURL		string
	ViewCount		int
	TextualUploadDate	string
	UploadDate		int
	isUploadDateApprox	int

}

func (r *DatabaseRepository) GetStream(streamURL string) Stream {
	var stream Stream

	err := r.Db.QueryRow("SELECT * FROM streams WHERE url = '" + streamURL  + "'").Scan(&stream.UID, &stream.ServiceID, &stream.URL, &stream.Title, &stream.StreamType, &stream.Duration, &stream.Uploader, &stream.UploaderURL, &stream.ThumbnailURL, &stream.ViewCount, &stream.TextualUploadDate, &stream.UploadDate, &stream.IsUploadDateApprox)
	if err != nil {
		panic(err)
	}

	return stream
}

func (r *DatabaseRepository) GetMaxPlaylistDisplayIndex() int {
	var maxDisplayIndex int
	err := r.Db.QueryRow("SELECT MAX(display_index) from playlists").Scan(&maxDisplayIndex)
	if err != nil {
		panic(err)
	}
	return maxDisplayIndex
}


func (r *DatabaseRepository) GetMaxPlaylistUID() int {
	var maxUID int
	err := r.Db.QueryRow("SELECT MAX(uid) from playlists").Scan(&maxUID)
	if err != nil {
		panic(err)
	}
	return maxUID
}

func (r *DatabaseRepository) GetPlaylistIDByName(playlistName string) int {
	var uid int
	err := r.Db.QueryRow("SELECT uid FROM playlists WHERE name = '" + playlistName + "'").Scan(&uid)
	if err != nil {
		panic(err)
	}
	return uid

}

func (r *DatabaseRepository) GetMaxPlaylistJoinIndex(playlistID int) int {
	var joinIndex *int
	err := r.Db.QueryRow("SELECT MAX(join_index) FROM playlist_stream_join WHERE playlist_id = " + strconv.Itoa(playlistID)).Scan(&joinIndex)
	if err != nil {
		panic(err)
	}

	//fmt.Println(*joinIndex)
	if (joinIndex == nil){
		return 0
	} else {
		return *joinIndex
	}
}


func (r *DatabaseRepository) GetMaxStreamUID() int {
	var maxUID int
	err := r.Db.QueryRow("SELECT MAX(uid) from streams").Scan(&maxUID)
	if err != nil {
		panic(err)
	}
	return maxUID
}

func (r *DatabaseRepository) CheckPlaylistExists(playlistName string) bool {
	var exists bool
	err := r.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM playlists WHERE name = '" + playlistName + "')").Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

func (r *DatabaseRepository) CheckStreamExists(vidURL string) bool {
	var exists bool
	err := r.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM streams WHERE url = '" + vidURL + "')").Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

func (r *DatabaseRepository) CheckStreamIsInPlaylist(stream Stream, playlistID int) bool {
	var exists bool
	err := r.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM playlist_stream_join WHERE playlist_id = " + strconv.Itoa(playlistID) + " AND stream_id = " + strconv.Itoa(stream.UID) + ")").Scan(&exists)

	if err != nil {
		panic(err)
	}
	return exists
}

func (r *DatabaseRepository) InsertPlaylist(playlist Playlist) {
	_, err := r.Db.Exec("INSERT INTO playlists (uid, name, is_thumbnail_permanent, thumbnail_stream_id, display_index) VALUES (?, ?, ?, ?, ?)", playlist.UID, playlist.Name, playlist.IsThumbnailPermanent, playlist.ThumbnailStreamID, playlist.DisplayIndex)

	if err != nil {
		panic(err)
	}
}

func (r *DatabaseRepository) InsertStream(stream Stream){
	_, err := r.Db.Exec("INSERT INTO streams (uid, service_id, url, title, stream_type, duration, uploader, uploader_url, thumbnail_url, view_count, textual_upload_date, upload_date, is_upload_date_approximation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", stream.UID, stream.ServiceID, stream.URL, stream.Title, stream.StreamType, stream.Duration, stream.Uploader, stream.UploaderURL, stream.ThumbnailURL, stream.ViewCount, stream.TextualUploadDate, stream.UploadDate, stream.IsUploadDateApprox)
	if err != nil {
		panic(err)
	}
}

func (r *DatabaseRepository) InsertStreamIntoPlaylist(stream Stream, playlistName string){
	
	playlistID := r.GetPlaylistIDByName(playlistName)
	if (r.CheckStreamIsInPlaylist(stream, playlistID)){
		return
	}


	_, err := r.Db.Exec("INSERT INTO playlist_stream_join (playlist_id, stream_id, join_index) VALUES (?, ?, ?)", playlistID, stream.UID, r.GetMaxPlaylistJoinIndex(playlistID) + 1)
	if err != nil{
		panic(err)
	}

	fmt.Printf("Added Video [%s] to Playlist [%s]\n", stream.Title, playlistName)
}

