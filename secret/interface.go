package secret

type CqpCall interface {
	IsGroupAdmin(group, qq uint64) bool
	IsGroupMaster(group, qq uint64) bool
}

var cqpCall CqpCall

func setCqpCall(v CqpCall) {
	cqpCall = v
}

func (b *Bot) isGroupAdmin(fromQQ uint64) bool {
	return cqpCall.IsGroupAdmin(b.Group, fromQQ)
}

func (b *Bot) isGroupMaster(fromQQ uint64) bool {
	return cqpCall.IsGroupMaster(b.Group, fromQQ)
}
