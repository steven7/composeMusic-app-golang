package models

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"strconv"
	"strings"

	//"github.com/aws/aws-sdk-go-v2/aws/session"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/steven7/go-createmusic/go/config"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FileType int

const (
	FileTypeImage FileType = 0
	FileTypeMusic FileType = 1
)

var (
	AWS_S3_BUCKET = "" // Bucket
)

var awsS3Client *s3.Client

func s3Client(c config.Config) *s3.Client {

	// Get values from .config file
	aws_access_key := c.AWS.AWSAccessKeyID
	aws_secret_key := c.AWS.AWSecretAccessKey
	aws_s3_region := c.AWS.AWSS3Region
	AWS_S3_BUCKET = c.AWS.S3.AWSS3Bucket

	// Create aws config
	creds := credentials.NewStaticCredentialsProvider(aws_access_key, aws_secret_key, "")
	credsConfig := awsConfig.WithCredentialsProvider(creds)
	regionConfig := awsConfig.WithRegion(aws_s3_region)
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), credsConfig, regionConfig)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	// create S3 client object
	s3Client := s3.NewFromConfig(cfg)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	return s3Client
}

type FileService interface {
	Create(trackID uint, fileReader io.Reader, filename string, ft FileType, c config.Config) error
	ByTrackID(trackID uint, ft FileType, filename string, c config.Config) (File, error)
	//ListByTrackID(trackID uint) ([]File, error)
	Delete(mf *File, ft FileType) error
}

func NewFileService() FileService {
	return &fileService{}
}

type fileService struct{}

func (fs *fileService) Create(trackID uint, fileReader io.Reader, filename string, ft FileType, c config.Config) error {
	if c.GetChannel() == "preprod" || c.GetChannel() == "prod" {
		return fs.CreateWithAWSS3(trackID, fileReader, filename, ft, c)
	} else {
		return fs.CreateLocal(trackID, fileReader, filename, ft)
	}
}

func (fs *fileService) CreateWithAWSS3(trackID uint, fileReader io.Reader, filename string, ft FileType, c config.Config) error {
	if awsS3Client == nil {
		awsS3Client = s3Client(c)
	}

	// Create file key.
	trackIDStr := strconv.FormatUint(uint64(trackID), 10)
	ftype := ""
	if ft == FileTypeMusic {
		ftype = "Music"
	} else {
		ftype = "Image"
	}
	fileName := trackIDStr + "/" + ftype + "/" + filename

	// Upload the file to AWS S3
	uploader := manager.NewUploader(awsS3Client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(fileName),
		Body:   fileReader,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	return nil
}

func (fs *fileService) CreateLocal(trackID uint, fileReader io.Reader, filename string, ft FileType) error {
	path, err := fs.mkDir(trackID, ft)
	if err != nil {
		return err
	}
	fmt.Println(path)
	// Clear directory before creating new one. We only need one file at a time.
	fs.ClearFiles(path)
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	fmt.Println(dst)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to the destination file
	_, err = io.Copy(dst, fileReader)
	if err != nil {
		return err
	}
	return nil
}

func (fs *fileService) Update(trackID uint, fileReader io.Reader, filename string, ft FileType, c config.Config) error {
	fmt.Println(" create music file")
	if c.GetChannel() == "preprod" || c.GetChannel() == "prod" {
		return fs.UpdateWithAWSS3(trackID, fileReader, filename, ft, c)
	} else {
		return fs.UpdateLocal(trackID, fileReader, filename, ft)
	}
}
func (fs *fileService) UpdateWithAWSS3(trackID uint, fileReader io.Reader, filename string, ft FileType, c config.Config) error {
	return nil
}

func (fs *fileService) UpdateLocal(trackID uint, fileReader io.Reader, filename string, ft FileType) error {
	return nil
}

/** private */
// Retuns open file. Make sure to close.
func (fs *fileService) CreateEmptyFileWithDirectory(trackID uint, filename string, ft FileType) (*os.File, error) {
	path, err := fs.mkDir(trackID, ft)
	if err != nil {
		return nil, err
	}

	fmt.Println(path)
	// Clear directory before creating new one. We only need one file at a time.
	fs.ClearFiles(path)

	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	fmt.Println(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (fs *fileService) ClearFiles(path string) {
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err == nil {
		fmt.Println(files)
		for _, imgStr := range files {
			//fmt.Println(imgStr)
			os.Remove(imgStr)
		}
	} else {
		//fmt.Println("lol file error")
		fmt.Println(err)
	}
}

func (fs *fileService) ByTrackID(trackID uint, ft FileType, filename string, c config.Config) (File, error) {
	if c.GetChannel() == "preprod" || c.GetChannel() == "prod" {
		return fs.ByTrackIDWithAWSS3(trackID, ft, filename, c)
	} else {
		return fs.ByTrackIDLocal(trackID, ft, c)
	}
}

func (fs *fileService) ByTrackIDWithAWSS3(trackID uint, ft FileType, filename string, c config.Config) (File, error) {
	if awsS3Client == nil {
		awsS3Client = s3Client(c)
	}

	ftype := ""
	if ft == FileTypeMusic {
		ftype = "Music"
	} else {
		ftype = "Image"
	}

	// Convert uint to string
	trackIDStr := strconv.FormatUint(uint64(trackID), 10)

	// Key to the file to be downloaded
	key := trackIDStr + "/" + ftype + "/" + filename

	// Name of the file where you want to save the downloaded file
	filenameL := "localfile_" + trackIDStr + "_" + ftype + "_" + filename

	// Create the file with directory
	newFile, err := fs.CreateEmptyFileWithDirectory(trackID, filenameL, ft)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer newFile.Close()

	downloader := manager.NewDownloader(awsS3Client)
	_, err = downloader.Download(context.TODO(), newFile, &s3.GetObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Get the file and parse out the path
	sfilename := strings.Split(filenameL, "/")
	localFilenameEnd := sfilename[len(sfilename)-1]

	file := File{
		Filename: localFilenameEnd,
		TrackID:  trackID,
	}

	return file, nil
}

func (fs *fileService) ByTrackIDLocal(trackID uint, ft FileType, c config.Config) (File, error) {
	path := fs.getDir(trackID, ft)
	// fmt.Println("musicfiles.go -  " + path)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	for _, s := range strings {
		fmt.Println(" files.go -  " + s)
	}
	//fmt.Println("images.go -  " + strings)
	if err != nil {
		return File{}, err
	}

	var filename string
	if len(strings) > 0 {
		filename = filepath.Base(strings[0])
	} else {
		filename = ""
	}

	ret := File{
		Filename: filename,
		TrackID:  trackID,
	}

	return ret, nil
}

func (fs *fileService) Delete(f *File, ft FileType) error {
	return os.Remove(f.RelativePath(ft))
}

func (fs *fileService) getDir(trackID uint, ft FileType) string {
	var dir string
	if ft == FileTypeMusic {
		dir = filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", trackID), "music")
	} else if ft == FileTypeImage {
		dir = filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", trackID), "cover")
	}
	return dir
}

func (fs *fileService) mkDir(galleryID uint, ft FileType) (string, error) {
	// filepath.Join will return a path like:
	//   images/galleries/123
	// We use filepath.Join instead of building the path
	// manually because the slashes and other characters
	// could vary between operating systems.
	var galleryPath string
	if ft == FileTypeMusic {
		galleryPath = filepath.Join("userfiles", "tracks",
			fmt.Sprintf("%v", galleryID), "music")
	} else if ft == FileTypeImage {
		galleryPath = filepath.Join("userfiles", "tracks",
			fmt.Sprintf("%v", galleryID), "cover")
	}
	// Create our directory (and any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

// File is used to represent images stored in a Gallery.
// File is NOT stored in the database, and instead
// references data stored on disk.
type File struct {
	TrackID  uint   `json:"trackId"`
	Filename string `json:"filename"`
}

// Path is used to build the absolute path used to reference this image
// via a web request.
//func (f *File) Path(ft FileType) string {
//	temp := url.URL{
//		Path: "/" + f.RelativePath(ft),
//	}
//	return temp.String()
//}

func (f *File) MusicPath() string {
	temp := url.URL{
		Path: f.RelativePath(FileTypeMusic),
	}
	return temp.String()
}

func (f *File) ImagePath() string {
	temp := url.URL{
		Path: f.RelativePath(FileTypeImage),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
func (f *File) RelativePath(ft FileType) string {
	// Convert the gallery ID to a string
	trackID := fmt.Sprintf("%v", f.TrackID)
	var path string
	if ft == FileTypeMusic {
		path = filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "music", f.Filename))
	} else if ft == FileTypeImage {
		path = filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "cover", f.Filename))
	}
	return path
}
