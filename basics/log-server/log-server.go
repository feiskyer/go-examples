package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

// LogServer serves log/exec api
type LogServer struct {
}

// Struct capturing information about an action ("GET", "POST", "DELETE", etc).
type action struct {
	Verb   string               // Verb identifying the action ("GET", "POST", "DELETE", etc).
	Path   string               // The path of the action
	Params []*restful.Parameter // List of parameters associated with the action.
}

func (s *LogServer) handleLogs(req *restful.Request, resp *restful.Response) {

}

func (s *LogServer) handleExec(req *restful.Request, resp *restful.Response) {

}

func (s *LogServer) registerContainerHandlers() error {
	ws := new(restful.WebService)

	nameParam := ws.PathParameter("id", "id of the container").DataType("string")
	params := []*restful.Parameter{nameParam}
	actions := []action{}

	actions = append(actions, action{"GET", "/container/{id}/logs", params})
	actions = append(actions, action{"POST", "/container/{id}/exec", params})

	for _, action := range actions {
		switch action.Verb {
		case "GET":
			doc := "read the specified containers"
			route := ws.GET(action.Path).To(s.handleLogs).
				Doc(doc).
				Operation("logs").
				Consumes(restful.MIME_XML, restful.MIME_JSON).
				Produces(restful.MIME_XML, restful.MIME_JSON)
			for _, param := range params {
				route.Param(param)
			}
			ws.Route(route)
			break
		case "POST":
			doc := "update the specified containers"
			route := ws.POST(action.Path).To(s.handleExec).
				Doc(doc).
				Operation("exec").
				Consumes(restful.MIME_XML, restful.MIME_JSON, "text/plain").
				Produces(restful.MIME_XML, restful.MIME_JSON)
			for _, param := range params {
				route.Param(param)
			}
			ws.Route(route)
			break
		default:
			return fmt.Errorf("unsupported action")
		}
	}

	restful.Add(ws)

	return nil
}

func main() {
	s := &LogServer{}
	s.registerContainerHandlers()
	http.ListenAndServe(":7777", nil)
}
