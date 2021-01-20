package lsp

import (
	"github.com/Mrs4s/MiraiGo/message"
	localdb "github.com/Sora233/Sora233-MiraiGo/lsp/buntdb"
	"github.com/Sora233/Sora233-MiraiGo/utils"
	"github.com/tidwall/buntdb"
	"strings"
	"time"
)

type KeySet struct{}

func (KeySet) GroupMessageImageKey(keys ...interface{}) string {
	return localdb.GroupMessageImageKey(keys...)
}

func (KeySet) GroupMuteKey(keys ...interface{}) string {
	return localdb.GroupMuteKey(keys...)
}

type StateManager struct {
	KeySet
}

func (s *StateManager) SaveMessageImageUrl(groupCode int64, messageID int32, msgs []message.IMessageElement) error {
	db, err := localdb.GetClient()
	if err != nil {
		return err
	}
	imgs := utils.MessageFilter(msgs, func(e message.IMessageElement) bool {
		return e.Type() == message.Image
	})
	var urls []string
	for _, img := range imgs {
		switch i := img.(type) {
		case *message.ImageElement:
			if i.Url != "" {
				urls = append(urls, i.Url)
			}
		case *message.GroupImageElement:
			if i.Url != "" {
				urls = append(urls, i.Url)
			}
		case *message.FriendImageElement:
			if i.Url != "" {
				urls = append(urls, i.Url)
			}
		}
	}
	if len(urls) > 0 {
		logger.WithField("group_code", groupCode).
			WithField("message_id", messageID).
			WithField("urls", urls).Debug("save image")
	} else {
		return nil
	}
	return db.Update(func(tx *buntdb.Tx) error {
		key := s.GroupMessageImageKey(groupCode, messageID)
		_, _, err := tx.Set(key, strings.Join(urls, " "), &buntdb.SetOptions{Expires: true, TTL: time.Minute * 30})
		return err
	})
}

func (s *StateManager) GetMessageImageUrl(groupCode int64, messageID int32) []string {
	db, err := localdb.GetClient()
	if err != nil {
		return nil
	}
	var result []string
	_ = db.View(func(tx *buntdb.Tx) error {
		key := s.GroupMessageImageKey(groupCode, messageID)
		val, err := tx.Get(key)
		if err == nil {
			result = strings.Split(val, " ")
		}
		return err
	})
	return result
}

func (s *StateManager) Muted(groupCode int64, uin int64, t int32) error {
	db, err := localdb.GetClient()
	if err != nil {
		return nil
	}
	return db.Update(func(tx *buntdb.Tx) error {
		key := s.GroupMuteKey(groupCode, uin)
		if t == 0 {
			_, err := tx.Delete(key)
			return err
		} else {
			_, _, err := tx.Set(key, "", &buntdb.SetOptions{
				Expires: true,
				TTL:     time.Second * time.Duration(t),
			})
			return err
		}
	})
}

func (s *StateManager) IsMuted(groupCode int64, uin int64) bool {
	db, err := localdb.GetClient()
	if err != nil {
		return false
	}
	var result = true
	db.View(func(tx *buntdb.Tx) error {
		key := s.GroupMuteKey(groupCode, uin)
		_, err := tx.Get(key)
		if err == buntdb.ErrNotFound {
			result = false
			return nil
		} else {
			return err
		}
	})
	return result
}

func (s *StateManager) FreshIndex() {
	db, _ := localdb.GetClient()
	db.CreateIndex(s.GroupMessageImageKey(), s.GroupMessageImageKey("*"), buntdb.IndexString)
}

func NewStateManager() *StateManager {
	return &StateManager{
		KeySet{},
	}
}