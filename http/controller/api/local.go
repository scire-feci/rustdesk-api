package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"net/http"
	"time"
)

// Local serves read-only status to processes on this host (see middleware.LocalOnly),
// e.g. Home Assistant polling device presence and the latest connection.
type Local struct {
}

// Heartbeats refresh last_online_time at most every 30s, so 90s tolerates two misses.
const localOnlineWindow = 90

type localDevice struct {
	Id         string `json:"id"`
	Alias      string `json:"alias"`
	Hostname   string `json:"hostname"`
	Os         string `json:"os"`
	Version    string `json:"version"`
	Online     bool   `json:"online"`
	LastOnline int64  `json:"last_online"`
	LastIp     string `json:"last_ip"`
}

type localConn struct {
	Id       uint   `json:"id"`
	Device   string `json:"device"`
	FromId   string `json:"from_id"`
	FromName string `json:"from_name"`
	Ip       string `json:"ip"`
	Type     int    `json:"type"`
	Time     int64  `json:"time"`
	Closed   bool   `json:"closed"`
}

// Status lists every device with its online state, plus the most recent connection.
func (l *Local) Status(c *gin.Context) {
	now := time.Now().Unix()
	devices := []localDevice{}
	online := 0
	for _, p := range service.AllService.PeerService.ListAll() {
		d := localDevice{
			Id:         p.Id,
			Alias:      p.Alias,
			Hostname:   p.Hostname,
			Os:         p.Os,
			Version:    p.Version,
			Online:     now-p.LastOnlineTime <= localOnlineWindow,
			LastOnline: p.LastOnlineTime,
			LastIp:     p.LastOnlineIp,
		}
		if d.Online {
			online++
		}
		devices = append(devices, d)
	}
	var last *localConn
	if ac := service.AllService.AuditService.LatestConn(); ac != nil {
		last = &localConn{
			Id:       ac.Id,
			Device:   ac.PeerId,
			FromId:   ac.FromPeer,
			FromName: ac.FromName,
			Ip:       ac.Ip,
			Type:     ac.Type,
			Time:     time.Time(ac.CreatedAt).Unix(),
			Closed:   ac.CloseTime != 0,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"online_window":   localOnlineWindow,
		"online_count":    online,
		"devices":         devices,
		"last_connection": last,
	})
}
