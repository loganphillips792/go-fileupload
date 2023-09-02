package api

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func (handler *Handler) CheckIfFileTypeIsSupported(file multipart.File) (bool, error) {
	// Only the first 512 bytes are used to sniff the content type
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return false, err
	}
	contentType := http.DetectContentType(buffer)
	handler.Logger.Info("CONTENT TYPE IS ", contentType)

	switch contentType {
	case "image/jpeg":
		return true, nil
	case "image/png":
		return true, nil
	default:
		return false, errors.New("file type not supported")
	}
}

// When user successfully uploads image, they can click "convert to black and white". The new image will
// show as a thumbnail, and then they can click download to download the new image
// https://stackoverflow.com/questions/42516203/converting-rgba-image-to-grayscale-golang
func (handler *Handler) ChangeImageToBlackAndWhite(filePath string) {
	handler.Logger.Infow("Hello", "msg", "Converting image to black and white...", "file path", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		handler.Logger.Error("Error ")
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		handler.Logger.Error(err)
	}

	// Create a new image with the same dimensions as the original
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// https://stackoverflow.com/a/42518487
	// Iterate through each pixel in the original image
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			oldPixel := img.At(x, y)
			pixel := color.GrayModel.Convert(oldPixel)
			newImg.Set(x, y, pixel)
		}
	}

	// Create a new file for the black and white image
	newFile, err := os.Create("uploads/bw.jpeg")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Encode the new image as a JPEG
	imageEncodeError := jpeg.Encode(newFile, newImg, nil)

	if imageEncodeError != nil {
		handler.Logger.Error("Error when encoding image: ", imageEncodeError)
	}

	time.Sleep(10 * time.Second)

	handler.Logger.Info("changeImageToBlackAndWhite go routine finished")
}
