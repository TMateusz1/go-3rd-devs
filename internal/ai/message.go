package ai

type Role string

const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

func UserMessage(content string) Message {
	return Message{
		Role:    User,
		Content: content,
	}
}

func SystemMessage(content string) Message {
	return Message{
		Role:    System,
		Content: content,
	}
}

func AssistantMessage(content string) Message {
	return Message{
		Role:    Assistant,
		Content: content,
	}
}
