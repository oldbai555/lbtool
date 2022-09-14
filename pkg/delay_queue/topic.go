package delay_queue

type Topic string

func (o Topic) String() string {
	return string(o)
}

func topicsToStrings(topics []Topic) []string {
	strings := make([]string, 0)
	for _, topic := range topics {
		strings = append(strings, string(topic))
	}
	return strings
}
