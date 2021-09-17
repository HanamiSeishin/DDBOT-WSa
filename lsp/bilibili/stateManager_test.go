package bilibili

import (
	"fmt"
	localdb "github.com/Sora233/DDBOT/lsp/buntdb"
	"github.com/Sora233/DDBOT/lsp/test"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
	"testing"
	"time"
)

func initStateManager(t *testing.T) *StateManager {
	sm := NewStateManager()
	assert.NotNil(t, sm)
	sm.FreshIndex(test.G1, test.G2)
	assert.Nil(t, sm.Start())
	return sm
}

func TestNewStateManager(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	sm := initStateManager(t)
	assert.NotNil(t, sm)
}

func TestStateManager_GetUserInfo(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)
	origUserInfo := NewUserInfo(test.UID1, test.ROOMID1, test.NAME1, "")
	assert.NotNil(t, origUserInfo)
	err := c.AddUserInfo(origUserInfo)
	assert.Nil(t, err)

	userInfo, err := c.GetUserInfo(test.UID1)
	assert.EqualValues(t, origUserInfo, userInfo)

	assert.NotNil(t, c.AddUserInfo(nil))
}

func TestStateManager_GetLiveInfo(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	origUserInfo := NewUserInfo(test.UID1, test.ROOMID1, test.NAME1, "")
	origLiveInfo := NewLiveInfo(origUserInfo, "", "", LiveStatus_Living)
	assert.NotNil(t, origLiveInfo)

	err := c.AddLiveInfo(origLiveInfo)
	assert.Nil(t, err)

	userInfo, err := c.GetUserInfo(test.UID1)
	assert.Nil(t, err)
	assert.NotNil(t, userInfo)
	assert.EqualValues(t, origUserInfo, userInfo)

	liveInfo, err := c.GetLiveInfo(test.UID1)
	assert.Nil(t, err)
	assert.NotNil(t, liveInfo)
	assert.EqualValues(t, origLiveInfo, liveInfo)

	liveInfo, err = c.GetLiveInfo(test.UID2)
	assert.Equal(t, buntdb.ErrNotFound, err)
	assert.Nil(t, liveInfo)

	err = c.DeleteLiveInfo(test.UID1)
	assert.Nil(t, err)

	liveInfo, err = c.GetLiveInfo(test.UID1)
	assert.Equal(t, buntdb.ErrNotFound, err)
	assert.Nil(t, liveInfo)

	assert.NotNil(t, c.AddLiveInfo(nil))
}

func TestStateManager_GetNewsInfo(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	origUserInfo := NewUserInfo(test.UID1, test.ROOMID1, test.NAME1, "")
	origNewsInfo := NewNewsInfo(origUserInfo, test.DynamicID1, test.TIMESTAMP1)

	err := c.AddNewsInfo(origNewsInfo)
	assert.Nil(t, err)

	userInfo, err := c.GetUserInfo(test.UID1)
	assert.Nil(t, err)
	assert.NotNil(t, userInfo)
	assert.EqualValues(t, origUserInfo, userInfo)

	newsInfo, err := c.GetNewsInfo(test.UID1)
	assert.Nil(t, err)
	assert.NotNil(t, newsInfo)
	assert.EqualValues(t, newsInfo, origNewsInfo)

	newsInfo, err = c.GetNewsInfo(test.UID2)
	assert.Equal(t, buntdb.ErrNotFound, err)
	assert.Nil(t, newsInfo)

	err = c.DeleteNewsInfo(test.UID1)
	assert.Nil(t, err)

	newsInfo, err = c.GetNewsInfo(test.UID1)
	assert.Equal(t, buntdb.ErrNotFound, err)
	assert.Nil(t, newsInfo)

	assert.NotNil(t, c.AddNewsInfo(nil))
}

func TestStateManager_DeleteNewsAndLiveInfo(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	origUserInfo := NewUserInfo(test.UID1, test.ROOMID1, test.NAME1, "")
	origNewsInfo := NewNewsInfo(origUserInfo, test.DynamicID1, test.TIMESTAMP1)
	origLiveInfo := NewLiveInfo(origUserInfo, "", "", LiveStatus_Living)
	assert.NotNil(t, origNewsInfo)
	assert.NotNil(t, origLiveInfo)

	assert.Nil(t, c.AddLiveInfo(origLiveInfo))
	assert.Nil(t, c.AddNewsInfo(origNewsInfo))

	assert.Nil(t, c.DeleteNewsAndLiveInfo(test.UID1))

	liveInfo, err := c.GetLiveInfo(test.UID1)
	assert.Nil(t, liveInfo)
	assert.NotNil(t, err)
	newsInfo, err := c.GetNewsInfo(test.UID1)
	assert.Nil(t, newsInfo)
	assert.NotNil(t, err)
}

func TestStateManager_CheckDynamicId(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	assert.True(t, c.CheckDynamicId(test.DynamicID1))

	replaced, err := c.MarkDynamicId(test.DynamicID1)
	assert.Nil(t, err)
	assert.False(t, replaced)

	assert.False(t, c.CheckDynamicId(test.DynamicID1))

	replaced, err = c.MarkDynamicId(test.DynamicID1)
	assert.Nil(t, err)
	assert.True(t, replaced)
}

func TestStateManager_IncNotLiveCount(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	assert.EqualValues(t, 1, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 2, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 3, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 4, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 5, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 6, c.IncNotLiveCount(test.UID1))

	assert.Nil(t, c.ClearNotLiveCount(test.UID1))
	assert.EqualValues(t, 1, c.IncNotLiveCount(test.UID1))
	assert.EqualValues(t, 2, c.IncNotLiveCount(test.UID1))
}

func TestStateManager_SetUidFirstTimestampIfNotExist(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	_, err := c.GetUidFirstTimestamp(test.UID2)
	assert.Equal(t, buntdb.ErrNotFound, err)

	assert.Nil(t, c.SetUidFirstTimestampIfNotExist(test.UID1, test.TIMESTAMP1))

	ts1, err := c.GetUidFirstTimestamp(test.UID1)
	assert.Nil(t, err)
	assert.Equal(t, test.TIMESTAMP1, ts1)

	assert.Nil(t, c.SetUidFirstTimestampIfNotExist(test.UID1, test.TIMESTAMP2))
	ts1, err = c.GetUidFirstTimestamp(test.UID1)
	assert.Nil(t, err)
	assert.Equal(t, test.TIMESTAMP1, ts1)

	assert.Nil(t, c.UnsetUidFirstTimestamp(test.UID1))

	ts1, err = c.GetUidFirstTimestamp(test.UID1)
	assert.Equal(t, buntdb.ErrNotFound, err)
}

func TestStateManager_ClearByMid(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	origUserInfo := NewUserInfo(test.UID1, test.ROOMID1, test.NAME1, "")
	origNewsInfo := NewNewsInfo(origUserInfo, test.DynamicID1, test.TIMESTAMP1)
	origLiveInfo := NewLiveInfo(origUserInfo, "", "", LiveStatus_Living)
	assert.NotNil(t, origNewsInfo)
	assert.NotNil(t, origLiveInfo)

	assert.Nil(t, c.AddLiveInfo(origLiveInfo))
	assert.Nil(t, c.AddNewsInfo(origNewsInfo))
	assert.Nil(t, c.SetUidFirstTimestampIfNotExist(test.UID1, test.TIMESTAMP1))
	assert.EqualValues(t, 1, c.IncNotLiveCount(test.UID1))

	assert.Nil(t, c.ClearByMid(test.UID1))

	userInfo, err := c.GetUserInfo(test.UID1)
	assert.NotNil(t, err)
	assert.Nil(t, userInfo)

	newsInfo, err := c.GetNewsInfo(test.UID1)
	assert.NotNil(t, err)
	assert.Nil(t, newsInfo)

	liveInfo, err := c.GetLiveInfo(test.UID1)
	assert.NotNil(t, err)
	assert.Nil(t, liveInfo)

	_, err = c.GetUidFirstTimestamp(test.UID1)
	assert.NotNil(t, err)
	assert.EqualValues(t, 1, c.IncNotLiveCount(test.UID1))

}

func TestGetCookieInfo(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	_ = initStateManager(t)

	cookieInfo, err := GetCookieInfo(test.NAME1)
	assert.EqualValues(t, buntdb.ErrNotFound, err)
	assert.Nil(t, cookieInfo)

	err = SetCookieInfo(test.NAME1, &LoginResponse_Data_CookieInfo{
		Cookies: []*LoginResponse_Data_CookieInfo_Cookie{
			{
				Name:  "name1",
				Value: "value1",
			},
			{
				Name:  "name2",
				Value: "value2",
			},
		},
		Domains: []string{"1"},
	})
	assert.Nil(t, err)

	cookieInfo, err = GetCookieInfo(test.NAME1)
	assert.Nil(t, err)
	assert.Len(t, cookieInfo.GetCookies(), 2)
	for idx, cookie := range cookieInfo.GetCookies() {
		assert.EqualValues(t, fmt.Sprintf("name%v", idx+1), cookie.GetName())
		assert.EqualValues(t, fmt.Sprintf("value%v", idx+1), cookie.GetValue())
	}
	assert.Len(t, cookieInfo.GetDomains(), 1)
	assert.Equal(t, "1", cookieInfo.GetDomains()[0])

	_, err = GetCookieInfo(test.NAME2)
	assert.NotNil(t, err)

	err = SetCookieInfo(test.NAME2, nil)
	assert.NotNil(t, err)
}

func TestStateManager_GetUserStat(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	userStat, err := c.GetUserStat(test.UID1)
	assert.NotNil(t, err)

	assert.NotNil(t, c.AddUserStat(nil, nil))

	userStat = NewUserStat(test.UID1, 1, 2)

	assert.Nil(t, c.AddUserStat(userStat, localdb.ExpireOption(time.Hour)))

	userStat = NewUserStat(test.UID2, 3, 4)

	assert.Nil(t, c.AddUserStat(userStat, localdb.ExpireOption(time.Hour)))

	userStat, err = c.GetUserStat(test.UID1)

	assert.EqualValues(t, 1, userStat.Following)
	assert.EqualValues(t, 2, userStat.Follower)
	assert.EqualValues(t, test.UID1, userStat.Mid)
}

func TestStateManager_SetGroupVideoOriginMarkIfNotExist(t *testing.T) {
	test.InitBuntdb(t)
	defer test.CloseBuntdb(t)

	c := initStateManager(t)

	assert.Nil(t, c.SetGroupVideoOriginMarkIfNotExist(test.G1, test.BVID1))
	assert.NotNil(t, c.SetGroupVideoOriginMarkIfNotExist(test.G1, test.BVID1))
}
