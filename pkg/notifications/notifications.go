package notifications

type notification struct {
	topic string
	msg   string
}

func newNotification(topic, msg string) notification {
	return notification{
		topic: topic,
		msg:   msg,
	}
}
