package flash

import (
	"app/pkg/cookie"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	MessageTypeSuccess = "success"
	MessageTypeWarning = "warning"
	MessageTypeDanger  = "danger"
)

type Message struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Valid    bool
	Duration int
}

const (
	flashMessageCookieName = "flash-message"
)

func Success(w http.ResponseWriter, msg string, duration int) error {
	return setMessage(w, Message{
		Type:     MessageTypeSuccess,
		Message:  msg,
		Valid:    true,
		Duration: duration,
	})
}

func Warning(w http.ResponseWriter, msg string, duration int) error {
	return setMessage(w, Message{
		Type:    MessageTypeWarning,
		Message: msg,
		Valid:   true,
	})
}
func Danger(w http.ResponseWriter, msg string, duration int) error {
	return setMessage(w, Message{
		Type:    MessageTypeDanger,
		Message: msg,
		Valid:   true,
	})
}

func setMessage(w http.ResponseWriter, flashMessage Message) error {
	b, err := json.Marshal(flashMessage)
	if err != nil {
		return fmt.Errorf("error marchalling the flash message which need to be saved in cookie:%w", err)
	}

	cookie.SetForOneYear(w, flashMessageCookieName, b64Encode(b))

	return nil
}

func Get(w http.ResponseWriter, r *http.Request) (Message, error) {
	m := Message{}
	c := cookie.Get(r, flashMessageCookieName)
	if c == "" {
		return m, nil
	}
	value, err := b64Decode(c)
	if err != nil {
		return m, fmt.Errorf("cannot decoding flash message : %w", err)
	}
	err = json.Unmarshal(value, &m)
	if err != nil {
		return m, fmt.Errorf("cannot unmarshaling flash message : %w", err)
	}
	cookie.Invalidate(w, flashMessageCookieName)
	return m, nil
}

func b64Encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func b64Decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
