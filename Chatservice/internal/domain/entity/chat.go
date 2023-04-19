package entity

import (
	"errors"

	"github.com/google/uuid"
)

type ChatConfig struct{
	Model *Model
	Temperature float32
	TopP float32
	N int
	Stop []string
	MaxTokens int
	PresencePenalty float32
	FrequencyPenalty float32
}

type Chat struct{
	ID string
	UserID string
	InitialSystemMessage *Message
	Messages []*Message
	EraseMessages []*Message
	Status string
	TokenUsage int
	Config *ChatConfig
}

func NewChat(userID string, initialSystemMenssage *Message, chatConfig *ChatConfig) (*Chat, error){
	chat := &Chat{
		ID: uuid.New().String(),
		UserID: userID,
		InitialSystemMessage: initialSystemMenssage,
		Status: "active",
		Config: chatConfig,
		TokenUsage: 0,
	}
	chat.AddMessege(initialSystemMenssage)

	if err := chat.Validate(); err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *Chat) Validate() error {
	if c.UserID == "" {
		return errors.New("user id is empty")
	}
	if c.Status != "active" && c.Status != "ended" {
		return errors.New("invalid status")
	}
	if c.Config.Temperature < 0 || c.Config.Temperature > 2 {
		return errors.New("invalid temperature")
	}
	// ... more validation for config
	return nil
}

 func (c *Chat) AddMessege(m *Message) error{
	if c.Status == "ended"{
		return errors.New("chat is ented, no more messages allowed")
	}
	for {
		if c.Config.Model.GetMaxTokens() >= m.GetQtdTokens()+c.TokenUsage{
			c.Messages = append(c.Messages, m)
			c.RefreshTokenUsage()
			break
		}
		c.EraseMessages = append(c.EraseMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokenUsage()
	}
	return nil
 }

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessages() int{
	return len(c.Messages)
}

func (c *Chat) Edn(){
	c.Status = "ended"
}

 func (c *Chat) RefreshTokenUsage(){
	c.TokenUsage = 0
	for m := range c.Messages{
		c.TokenUsage += c.Messages[m].GetQtdTokens()
	}
 }