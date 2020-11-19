package rmupdate

import (
	"github.com/saleemrashid/rmupdate/omaha"
)

type RequestParams struct {
	AppID          string
	Group          string
	MachineType    string
	OSIdentifier   string
	OSVersion      string
	Platform       string
	ReleaseVersion string
	SerialNumber   string
}

func (p *RequestParams) Build() *omaha.Request {
	return &omaha.Request{
		Protocol:       "3.0",
		Version:        p.ReleaseVersion,
		RequestID:      NewUUID().String(),
		SessionID:      NewUUID().String(),
		UpdaterVersion: "0.4.2",
		InstallSource:  "ondemandupdate",
		IsMachine:      1,
		OS: omaha.OS{
			Version:     p.OSIdentifier + " " + p.OSVersion,
			Platform:    p.Platform,
			ServicePack: p.ReleaseVersion + "_" + p.MachineType,
			Arch:        p.MachineType,
		},
		Apps: []omaha.AppRequest{{
			ID:                   p.AppID,
			Version:              p.ReleaseVersion,
			Track:                p.Group,
			AdditionalParameters: p.Group,
			BootID:               NewUUID().String(),
			OEM:                  p.SerialNumber,
			OEMVersion:           p.OSVersion,
			AlephVersion:         p.ReleaseVersion,
			MachineID:            NewMachineID(),
			Lang:                 "en-US",
			DeltaOK:              false,
			Ping: omaha.PingRequest{
				Active: 1,
			},
			Events: []omaha.EventRequest{{
				Type:   omaha.EventTypeUpdateComplete,
				Result: omaha.EventResultSuccessReboot,
			}},
		}},
	}
}
