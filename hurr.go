package hurr

import (
	"bytes"
	"errors"
	"strings"
)

type TranslatedError interface {
	Error() string
}

type errorMessage struct {
	manager      *Manager
	template     string
	translations []TranslatedError
}

func (e *errorMessage) Set(language, message string) error {
	i, err := e.manager.languageIndex(language)
	if err != nil {
		return err
	}

	e.translations[i] = errors.New(message)

	return nil
}

func (e *errorMessage) findValues(text string) map[string]string {
	j := 0
	startVar := -1

	data := map[string]string{}
	startVal := -1

	for i := 0; i < len(e.template)-2; i++ {
		switch e.template[i : i+2] {
		case "{{":
			startVar = i + 2
			startVal = j
		case "}}":
			if startVar == -1 {
				continue
			}

			variable := strings.TrimSpace(e.template[startVar:i])
			var value string

			scanUntil := e.template[i+2]
			for k := j; k < len(text); k++ {
				if text[k] == scanUntil {
					value = text[startVal:k]
					j = k - 2
					break
				}
			}
			data[variable] = value
			startVar = -1
		}
		if startVar == -1 {
			j++
		}
	}

	return data
}

func (e *errorMessage) populate(index int, data map[string]string) string {
	var buf bytes.Buffer
	tmpl := e.translations[index].Error()

	startVar := -1
	prevVar := 0
	for i := 0; i < len(tmpl)-2; i++ {
		switch tmpl[i : i+2] {
		case "{{":
			startVar = i + 2
		case "}}":
			if startVar == -1 {
				continue
			}

			variable := strings.TrimSpace(tmpl[startVar:i])
			buf.WriteString(tmpl[prevVar : startVar-2])
			buf.WriteString(data[variable])
			prevVar = i + 2
		}
	}

	buf.WriteString(tmpl[prevVar:])
	return buf.String()
}

type Manager struct {
	Errors    []*errorMessage
	languages []string
}

func (m *Manager) languageIndex(code string) (int, error) {
	for i, langCode := range m.languages {
		if code == langCode {
			return i, nil
		}
	}

	return -1, errors.New("invalid language code")
}

func (m *Manager) Add(tmpl string) *errorMessage {
	err := errorMessage{
		manager:      m,
		template:     tmpl,
		translations: make([]TranslatedError, len(m.languages)),
	}
	m.Errors = append(m.Errors, &err)

	return &err
}

func (m *Manager) findErrorTemplate(text string) (*errorMessage, map[string]string) {
	i := 0
	var errMsg *errorMessage
	pos := make([]int, len(m.Errors))
	matchUntil := make([]byte, len(m.Errors))
	possibleMatches := len(m.Errors)
	for possibleMatches > 0 && errMsg == nil {
		for j, msg := range m.Errors {
			p := pos[j]

			if pos[j] == -1 {
				continue
			}

			if matchUntil[j] == text[i] {
				matchUntil[j] = 0
			}

			if msg.template[p] == text[i] {
				pos[j]++
			} else if msg.template[p] == '{' && msg.template[p+1] == '{' {
				for q := p; q < len(msg.template); q++ {
					if msg.template[q] == '}' && msg.template[q+1] == '}' {
						pos[j] = q + 2
						matchUntil[j] = msg.template[pos[j]]
						break
					}
				}
			} else if matchUntil[j] == 0 {
				pos[j] = -1
				possibleMatches -= 1
			}

			if pos[j] == len(msg.template) {
				errMsg = msg
				break
			}
		}
		i++
	}

	return errMsg, nil
}

func (m *Manager) Get(language string, err error) (string, error) {
	langIndex, ierr := m.languageIndex(language)
	if ierr != nil {
		return "", ierr
	}

	errMsg, _ := m.findErrorTemplate(err.Error())
	if errMsg == nil {
		return "", errors.New("no matching error messages found")
	}

	data := errMsg.findValues(err.Error())
	str := errMsg.populate(langIndex, data)

	return str, nil
}

func NewManager(languages []string) *Manager {
	mgr := Manager{
		[]*errorMessage{},
		languages,
	}

	return &mgr
}
