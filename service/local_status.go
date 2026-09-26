package service

import "github.com/lejianwen/rustdesk-api/v2/model"

// ListAll returns every device, most recently seen first.
func (ps *PeerService) ListAll() (peers []*model.Peer) {
	DB.Order("last_online_time desc").Find(&peers)
	return
}

// LatestConn returns the most recent connection audit, or nil if there is none.
func (as *AuditService) LatestConn() *model.AuditConn {
	ac := &model.AuditConn{}
	if err := DB.Order("id desc").First(ac).Error; err != nil {
		return nil
	}
	return ac
}
