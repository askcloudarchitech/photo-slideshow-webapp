package photoprocessing

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"runtime"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func DetectAndConvertImage(fileHeader *multipart.FileHeader) (image.Image, int, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, err
	}

	fileType, err := getFileType(file)
	if err != nil {
		return nil, 0, err
	}

	file.Seek(0, 0)

	var img []byte

	switch fileType {
	case "image/heic", "image/heif", "application/octet-stream":
		img, err = heicToJpeg(file)
		if err != nil {
			return nil, 0, err
		}
		break
	case "image/jpeg", "image/jpg":
		img, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, 0, err
		}
		break
	default:
		return nil, 0, fmt.Errorf("unsupported file type")
	}

	image, err := jpeg.Decode(bytes.NewBuffer(img))
	if err != nil {
		log.Println("could not convert to img data")
		return nil, 0, err
	}

	exifData, err := exif.Decode(bytes.NewBuffer(img))
	if err != nil {
		log.Println("could not extract exif data")
		return image, int(time.Now().Unix()), nil
	}

	date, err := exifData.DateTime()
	if err != nil {
		log.Printf("could not get datetime, returning timestamp" + err.Error())
		date = time.Now()
	}

	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		// tag not present
		log.Println("no orientation tag")
		return image, int(date.Unix()), nil
	}

	if tag.Count == 1 && tag.Format() == tiff.IntVal {
		orientation, err := tag.Int(0)
		if err != nil {
			return image, int(date.Unix()), nil
		}

		log.Println(orientation)

		switch orientation {
		case 3: // rotate 180
			image = imaging.Rotate180(image)
		case 6: // rotate 270
			image = imaging.Rotate270(image)
		case 8: //rotate 90
			image = imaging.Rotate90(image)
		}
	}

	return image, int(date.Unix()), nil
}

func getFileType(file multipart.File) (string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)

	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	return filetype, nil
}
