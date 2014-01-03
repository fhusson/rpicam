// +build linux

package rpicam

const (
	raspistillCommand          = "raspistill"
	raspistillDefaultArguments = "-hf -vf -o - -t 500 --nopreview"
)
