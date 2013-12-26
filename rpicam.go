package rpicam

import (
	"os/exec"
	"strings"
)

const (
	raspistillCommand          = "raspistill"
	raspistillDefaultArguments = "-hf -vf -o - -t 500 --nopreview"
)

type Manager struct {
	latest []byte
	error  []byte
	call   chan Request
}

type Response struct {
	Err  error
	Data []byte
}

type Request struct {
	answer chan Response
	args   string
}

func NewManager() *Manager {
	pm := Manager{}
	pm.call = make(chan Request)

	return &pm
}

func (pm *Manager) NewShot(args string) Response {
	req := Request{}
	req.answer = make(chan Response)
	req.args = args

	// we pass the order
	pm.call <- req

	// we send back the response
	return <-req.answer
}

func (pm *Manager) LatestShot() Response {
	var r Response
	if len(pm.latest) > 0 {
		r = Response{}
		r.Data = pm.latest
	} else {
		r = pm.NewShot("")
	}
	return r
}

func (pm *Manager) Serve() {
	for {
		select {
		case request := <-pm.call:
			response := Response{}

			if request.args != "latest" {
				args := raspistillDefaultArguments
				if len(request.args) > 0 {
					args = args + " " + request.args
				}
				argsArray := strings.Split(args, " ")
				out, err := exec.Command(raspistillCommand, argsArray...).CombinedOutput()
				response.Err = err
				response.Data = out
				pm.latest = out
			} else {
				response.Data = pm.latest
			}
			request.answer <- response
		}
	}
}
