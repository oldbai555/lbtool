package fightinglandlords

import "github.com/oldbai555/lbtool/pkg/lberr"

var (
	NotCanJoinDesk          = lberr.NewErr(1, "请先退出当前房间加入其他房间")
	DeskPlayerUpperLimit    = lberr.NewErr(2, "房间已经满人了,请选择其他房间")
	DeskPlayerNotUpperLimit = lberr.NewErr(3, "房间还未满人,请稍等")
	PlayerNotFoundLeave     = lberr.NewErr(4, "找不到要退房的用户")
	DeskNotFound            = lberr.NewErr(5, "找不到房间")
	DeskNotLandlordRound    = lberr.NewErr(6, "没人叫地主")
	DeskPlayerUnReady       = lberr.NewErr(7, "还有玩家未准备")
)
