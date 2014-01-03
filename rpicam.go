package rpicam

import (
	"errors"
	"os/exec"
	"strings"
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
				if err == nil {
					// we check the magic number FF D8
					if len(out) > 2 && out[0] == 255 && out[1] == 216 {
						response.Data = out
						pm.latest = out
					} else {
						response.Err = errors.New(string(out))
					}
				} else {
					response.Err = err
				}
			} else {
				response.Data = pm.latest
			}
			request.answer <- response
		}
	}
}
