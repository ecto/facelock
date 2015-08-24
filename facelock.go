package main

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/lazywei/go-opencv/opencv"
)

func main() {
	cap := opencv.NewCameraCapture(0)
	if cap == nil {
		panic("can not open camera")
	}
	defer cap.Release()

	_, currentfile, _, _ := runtime.Caller(0)
	cascade := opencv.LoadHaarClassifierCascade(path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml"))
	timeout := time.Second * 10
	lastFace := time.Now()

	for {
		if cap.GrabFrame() {
			img := cap.RetrieveFrame(1)
			if img == nil {
				continue
			}
			hasFaces, _ := ProcessImage(img, cascade)
			if hasFaces {
				lastFace = time.Now()
			} else {
				delta := time.Since(lastFace)
				fmt.Println("no faces for %ts", delta)
				if delta > timeout {
					fmt.Println("lock!")
					lock()
					time.Sleep(time.Second * 60)
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func ProcessImage(img *opencv.IplImage, cascade *opencv.HaarCascade) (bool, error) {
	faces := cascade.DetectObjects(img)
	hasFaces := len(faces) > 0
	return hasFaces, nil
}

func lock() {
	cmd := "/System/Library/Frameworks/ScreenSaver.framework/Resources/ScreenSaverEngine.app/Contents/MacOS/ScreenSaverEngine"
	out, err := exec.Command(cmd).Output()
	fmt.Println(out)
	fmt.Println(err)
}
