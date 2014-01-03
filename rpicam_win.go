// +build !linux

package rpicam

const (
	raspistillCommand          = "cmd"
	raspistillDefaultArguments = "/c echo raspistill -hf -vf -o - -t 500 --nopreview"
)
