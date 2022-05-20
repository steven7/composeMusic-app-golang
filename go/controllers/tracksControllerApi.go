package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/steven7/go-createmusic/go/config"
	"github.com/steven7/go-createmusic/go/models"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

func NewTracksAPI(ts models.TrackService, fs models.FileService, r *mux.Router, c config.Config) *TrackController {
	return &TrackController{
		ts: ts,
		fs: fs,
		//is:                	   is,
		//mfs: 				   mfs,
		r: r,
		c: c,
	}
}

// GET /api/tracks/index
func (t *TrackController) IndexWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("IndexWithAPI!!")

	//
	// validate jwt token is done with middleware
	//
	var info models.TrackIndexJson
	err := ParseJSONParameters(w, r, &info)
	if err != nil {
		fmt.Println("Could not parse parameters. Error response sent")
		fmt.Println("Error: ", err)
		return
	}

	uintUserID := info.UserID
	fmt.Println(uintUserID)

	tracks, err := t.ts.ByUserID(uint(uintUserID))
	fmt.Println("json")
	for _, track := range tracks {
		fmt.Println(track)
	}

	if err != nil {
		errorData := models.Error{
			Title:  "Error fetching tracks for user",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteJson(w, tracks)

}

// GET /api/tracks/one
func (t *TrackController) GetTrackWithAPI(w http.ResponseWriter, r *http.Request) {

	//
	// validate jwt token is done with middleware
	//
	var info models.OneTrackJson
	ParseJSONParameters(w, r, &info)
	uintTrackID := info.TrackID

	track, err := t.ts.ByID(uint(uintTrackID))
	// if error choose way to respond
	if err != nil {
		errorData := models.Error{
			Title:  "Error fetching one track",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	// music file
	//musicfile, err := t.fs.ByTrackID(track.ID, models.FileTypeMusic)
	//if err != nil {
	//	errorData := models.Error {
	//		Title:  "Error fetching the track file",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	// cover image
	//coverimage, err := t.fs.ByTrackID(track.ID, models.FileTypeImage)
	//if err != nil {
	//	errorData := models.Error {
	//		Title:  "Error fetching the cover image file",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	WriteJson(w, track)
}

// GET /api/tracks/one/coverimage
func (t *TrackController) GetTrackCoverFileWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetTrackCoverFileWithAPI")

	var info models.OneTrackJson
	err := ParseJSONParameters(w, r, &info)
	if err != nil {
		fmt.Println("Could not parse parameters. Error response sent")
		return
	}

	uintTrackID := info.TrackID
	filename := info.FileName

	fmt.Println(uintTrackID)

	// cover image
	coverimage, err := t.fs.ByTrackID(uint(uintTrackID), models.FileTypeImage, filename, t.c)
	if err != nil {
		errorData := models.Error{
			Title:  "Error fetching the cover image file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteFile(w, coverimage.ImagePath())

}

// GET /api/tracks/one/musicfile
func (t *TrackController) GetTrackMusicFileWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetTrackMusicFileWithAPI")

	var info models.OneTrackJson
	err := ParseJSONParameters(w, r, &info)
	if err != nil {
		fmt.Println("Could not parse parameters. Error response sent")
		return
	}

	uintTrackID := info.TrackID
	filename := info.FileName

	// music file

	musicfile, err := t.fs.ByTrackID(uint(uintTrackID), models.FileTypeMusic, filename, t.c)
	if err != nil {
		errorData := models.Error{
			Title:  "Error fetching the track file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteFile(w, musicfile.MusicPath())

}

// POST /api/tracks/createlocal
func (t *TrackController) CreateLocalWithAPI(w http.ResponseWriter, r *http.Request) {

	//
	// validate jwt token is done with middleware
	//

	r.ParseMultipartForm(256)

	if r.MultipartForm == nil {
		errorData := models.Error{
			Title:  "Error creating track",
			Detail: "Multipart/form data not correctly included in request.",
		}
		WriteJson(w, errorData)
		return
	}

	f := r.MultipartForm

	userID, err := strconv.ParseUint(f.Value["userID"][0], 10, 32)
	if err != nil {
		userID = 0
		fmt.Println(err)
	}
	title := f.Value["title"][0]
	artist := f.Value["artist"][0]
	desc := f.Value["desc"][0]

	fmt.Println("the multipart form %s", f)

	///
	///
	// Get filenames from headers first.
	///
	///

	//
	//
	// cover image
	//
	//

	imgArr := f.File["coverimage"]
	var imageHeader *multipart.FileHeader
	var coverImageFileName string = ""
	if f != nil && f.File["coverimage"] != nil && len(imgArr) > 0 {
		imageHeader = f.File["coverimage"][0]
		coverImageFileName = imageHeader.Filename
	}

	//
	//
	// music file
	//
	//

	if f == nil || f.File["musicfile"] == nil || len(f.File["musicfile"]) == 0 {
		errorData := models.Error{
			Title:  "Error creating track",
			Detail: "A valid music file must be included",
		}
		WriteJson(w, errorData)
		return
	}

	musicHeader := f.File["musicfile"][0]
	musicfile, err := musicHeader.Open()

	musicFileName := musicHeader.Filename

	// Create track object in memory. It will be stored once music file and cover get processed.
	track := models.Track{
		Title:              title, //trackData.Track.Title,
		Artist:             artist,
		Description:        desc,         //"Created with create local track api endpoint",
		UserID:             uint(userID), //trackData.UserID, //trackData.User.ID,
		CoverImageFilename: coverImageFileName,
		MusicFilename:      musicFileName,
	}

	//
	//
	// Create track file in database
	//
	//

	if err := t.ts.Create(&track); err != nil {
		// app.serverError(w, err)
		errorData := models.Error{
			Title:  "Error creating track",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	///
	///
	// Write files once the track object is created.
	///
	///

	//
	//
	// cover image
	//
	//

	//
	// Headers defined earlier above
	//
	if imageHeader != nil {

		imagefile, err := imageHeader.Open()
		if err != nil {
			errorData := models.Error{
				Title:  "Error with cover image",
				Detail: err.Error(),
			}
			WriteJson(w, errorData)
			return
		}
		defer imagefile.Close()

		err = t.fs.Create(track.ID, imagefile, imageHeader.Filename, models.FileTypeImage, t.c)
		if err != nil {
			errorData := models.Error{
				Title:  "Error uploading cover image",
				Detail: err.Error(),
			}
			WriteJson(w, errorData)
			return
		}
		track.CoverImageFilename = imageHeader.Filename

	}

	//
	//
	// music file
	//
	//

	err = t.fs.Create(track.ID, musicfile, musicFileName, models.FileTypeMusic, t.c)
	if err != nil {
		errorData := models.Error{
			Title:  "Error uploading file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	if err != nil {
		errorData := models.Error{
			Title:  "Error with uploaded file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}
	defer musicfile.Close()

	//
	//
	// Update track file with info connecting it to the icon and music file.
	//
	//

	if err := t.ts.Update(&track); err != nil {
		// app.serverError(w, err)
		errorData := models.Error{
			Title:  "Error creating track",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	///
	createLocalTrackResponse := models.CreateLocalTrackResponseJson{
		Success: true,
		Message: "Track successfully created!!",
		Track:   track,
	}

	WriteJson(w, createLocalTrackResponse)
}

// POST /tracks/createWithComposeAI
func (t *TrackController) CreateWithComposeAI(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(256)
	if err != nil {
		fmt.Println("Could not parse parameters. Error response sent")
		return
	}

	f := r.MultipartForm

	var userID string
	var title string
	var artist string
	var desc string
	var composeType string
	if len(f.Value["userID"]) > 0 {
		userID = f.Value["userID"][0]
	}
	if len(f.Value["title"]) > 0 {
		title = f.Value["title"][0]
	}
	if len(f.Value["artist"]) > 0 {
		artist = f.Value["artist"][0]
	}
	if len(f.Value["desc"]) > 0 {
		desc = f.Value["desc"][0]
	}
	if len(f.Value["compose_type"]) > 0 {
		composeType = f.Value["compose_change_type"][0]
	}

	///
	///

	type Dictionary map[string]interface{}

	jsonDict := Dictionary{
		"userID":       string(userID),
		"title":        title,
		"artist":       artist,
		"desc":         desc,
		"compose_type": composeType,
	}

	body, _ := json.Marshal(jsonDict)

	// Post to Compose machine learning service

	// Make channel
	composeChan := make(chan *http.Response)

	// Use go routine
	go models.SendPostAsync(body, composeChan)

	// Make request to compose ml service
	composeResponse := <-composeChan
	defer composeResponse.Body.Close()
	bytes, _ := ioutil.ReadAll(composeResponse.Body)
	fmt.Println(string(bytes))

	// response from compose service
	success := true
	responseTitle := "Track successfully composed by the machine learning ai dj."

	if composeResponse.StatusCode != 200 {
		success = false
		responseTitle = "Error with compose music service"
	}

	createLocalTrackResponse := models.CreateLocalTrackResponseJson{
		Success: success,
		Title:   responseTitle,
		Message: string(bytes),
	}

	fmt.Println("the response from golang " + responseTitle + " " + string(bytes))

	//WriteJson(w, createLocalTrackResponse)
	WriteJsonWithStatus(w, createLocalTrackResponse, 200)
}

// POST /track/create
func (t *TrackController) CreateDJWithAPI(w http.ResponseWriter, r *http.Request) {

	//userId := 2 // default text id
	//trackId := 10 // track 10 is 0x00ss.jpg cover image and Tchaikovsky-Waltz-op39-no8 file
}

// POST /track/create
func (t *TrackController) CreateDJ_Jazz_WithAPI(w http.ResponseWriter, r *http.Request) {

	//userId := 2 // default text id
	//trackId := 10 // track 10 is 0x00ss.jpg cover image and Tchaikovsky-Waltz-op39-no8 file

	//cmd := exec.Command("script.py")
	//
	//cmd := exec.Command("cmd", python, script)
	//out, err := cmd.Output()
	//fmt.Println(string(out))

}

// GET /api/tracks/update
func (t *TrackController) UpdateTrackWithAPI(w http.ResponseWriter, r *http.Request) {

}

// GET /api/tracks/delete
func (t *TrackController) DeleteTrackWithAPI(w http.ResponseWriter, r *http.Request) {

	var info models.OneTrackJson
	ParseJSONParameters(w, r, &info)
	uiTrackID := uint(info.TrackID)

	fmt.Println(uiTrackID)

	var additionalInfo = ""

	// Cover Image

	// Find cover image first before deleting it. Consider changing this later.
	coverimage, err := t.fs.ByTrackID(uiTrackID, models.FileTypeImage, "", t.c)
	// Delete cover image if found.
	if err = t.fs.Delete(&coverimage, models.FileTypeImage); err != nil {
		//errorData := models.Error {
		//	Title:  "Error deleting the cover image",
		//	Detail: err.Error(),
		//}
		//WriteJson(w, errorData)
		//return
		additionalInfo += " Could not find cover image file to delete."
	}

	// Music file

	// Find music file first before deleting it. Consider changing this later.
	musicfile, err := t.fs.ByTrackID(uiTrackID, models.FileTypeMusic, "", t.c)
	// Delete music file if found.
	if err = t.fs.Delete(&musicfile, models.FileTypeMusic); err != nil {
		//	errorData := models.Error {
		//		Title:  "Error deleting the track file",
		//		Detail: err.Error(),
		//	}
		//	WriteJson(w, errorData)
		//	return
		additionalInfo += " Could not find music file to delete."
	}

	if err := t.ts.Delete(uiTrackID); err != nil {
		errorData := models.Error{
			Title:  "Error deleting track metadata. " + additionalInfo,
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	deleteTrackResponse := models.SuccessJson{
		Success: true,
		Message: "Track successfully deleted." + additionalInfo,
	}

	WriteJson(w, deleteTrackResponse)
}

// GET /track/download
func (t *TrackController) DownloadFileWithAPI(w http.ResponseWriter, r *http.Request) {

}

// GET /track/stream
func (t *TrackController) StreamWithAPI(w http.ResponseWriter, r *http.Request) {

}
