package omaha

import (
	"encoding/xml"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeDownloadComplete
	EventTypeInstallComplete
	EventTypeUpdateComplete
	EventTypeUpdateDownloadStarted
	EventTypeUpdateDownloadFinished
)

type EventResult int

const (
	EventResultError EventResult = iota
	EventResultSuccess
	EventResultSuccessReboot
	EventResultUpdateDeferred
)

type Request struct {
	XMLName        xml.Name     `xml:"request"`
	Protocol       string       `xml:"protocol,attr"`
	Version        string       `xml:"version,attr"`
	RequestID      string       `xml:"requestid,attr"`
	SessionID      string       `xml:"sessionid,attr"`
	UpdaterVersion string       `xml:"updaterversion,attr"`
	InstallSource  string       `xml:"installsource,attr"`
	IsMachine      int          `xml:"ismachine,attr"`
	OS             OS           `xml:"os"`
	Apps           []AppRequest `xml:"app"`
}

type OS struct {
	XMLName     xml.Name `xml:"os"`
	Version     string   `xml:"version,attr"`
	Platform    string   `xml:"platform,attr"`
	ServicePack string   `xml:"sp,attr"`
	Arch        string   `xml:"arch,attr"`
}

type AppRequest struct {
	XMLName              xml.Name           `xml:"app"`
	ID                   string             `xml:"appid,attr"`
	Version              string             `xml:"version,attr"`
	Track                string             `xml:"track,attr"`
	AdditionalParameters string             `xml:"ap,attr"`
	BootID               string             `xml:"bootid,attr"`
	OEM                  string             `xml:"oem,attr"`
	OEMVersion           string             `xml:"oemversion,attr"`
	AlephVersion         string             `xml:"alephversion,attr"`
	MachineID            string             `xml:"machineid,attr"`
	Lang                 string             `xml:"lang,attr"`
	Board                string             `xml:"board,attr"`
	HardwareClass        string             `xml:"hardware_class,attr"`
	DeltaOK              bool               `xml:"delta_okay,attr"`
	NextVersion          string             `xml:"nextversion,attr"`
	Brand                string             `xml:"brand,attr"`
	Client               string             `xml:"client,attr"`
	Ping                 PingRequest        `xml:"ping"`
	UpdateCheck          UpdateCheckRequest `xml:"updatecheck"`
	Events               []EventRequest     `xml:"event"`
}

type PingRequest struct {
	XMLName xml.Name `xml:"ping"`
	Active  int      `xml:"active,attr"`
}

type UpdateCheckRequest struct {
	XMLName xml.Name `xml:"updatecheck"`
}

type EventRequest struct {
	XMLName         xml.Name    `xml:"event"`
	Type            EventType   `xml:"eventtype,attr"`
	Result          EventResult `xml:"eventresult,attr"`
	PreviousVersion string      `xml:"previousversion,attr"`
}

type Response struct {
	XMLName  xml.Name      `xml:"response"`
	Protocol string        `xml:"protocol,attr"`
	Apps     []AppResponse `xml:"app"`
}

type AppResponse struct {
	XMLName     xml.Name            `xml:"app"`
	ID          string              `xml:"appid,attr"`
	UpdateCheck UpdateCheckResponse `xml:"updatecheck"`
}

type UpdateCheckResponse struct {
	XMLName  xml.Name `xml:"updatecheck"`
	Status   string   `xml:"status,attr"`
	URLs     []URL    `xml:"urls>url"`
	Manifest Manifest `xml:"manifest"`
}

type URL struct {
	XMLName  xml.Name `xml:"url"`
	CodeBase string   `xml:"codebase,attr"`
}

type Manifest struct {
	XMLName  xml.Name  `xml:"manifest"`
	Version  string    `xml:"version,attr"`
	Packages []Package `xml:"packages>package"`
	Actions  []Action  `xml:"actions>action"`
}

type Package struct {
	XMLName xml.Name `xml:"package"`
	SHA1    Hash     `xml:"hash,attr"`
	Name    string   `xml:"name,attr"`
	Size    int64    `xml:"size,attr"`
}

type Action struct {
	XMLName xml.Name `xml:"action"`
	SHA256  Hash     `xml:"sha256,attr"`
	Event   string   `xml:"event,attr"`
}
