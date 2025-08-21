package messagetype

type Type string

const (
	Text    Type = "text"
	Photo   Type = "photo"
	Voice   Type = "voice"
	Sticker Type = "sticker"
	System  Type = "system"
)
