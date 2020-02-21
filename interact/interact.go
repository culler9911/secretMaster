package interact

import (
	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
)

func GetGroupMemberInfo(group, qq int64, noCatch bool) cqp.GroupMember {
	return cqp.GetGroupMemberInfo(group, qq, noCatch)
}
