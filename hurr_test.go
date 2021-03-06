package hurr

import (
	"errors"
	"fmt"
	"testing"
)

func TestLangIndex(t *testing.T) {
	mgr := NewManager([]string{"English", "Norwegian Bokmål"})
	if i, err := mgr.languageIndex("English"); i != 0 && err != nil {
		t.Fail()
	}
	if i, err := mgr.languageIndex("Norwegian Bokmål"); i != 0 && err != nil {
		t.Fail()
	}
}

func TestSimple(t *testing.T) {
	mgr := NewManager([]string{"English", "Norwegian Bokmål"})

	errMsg := mgr.Add(`pq: duplicate key value violates unique constraint "{{ table }}_{{ column }}_key"`)
	errMsg.Set("English", `This {{ column }} already exists in {{ table }}.`)
	errMsg.Set("Norwegian Bokmål", `Denne {{ column }} eksisterer allerede i {{ table }}.`)

	err := errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)
	if msg, _ := mgr.Get("English", err); msg != "This email already exists in users." {
		t.Fail()
	}

	if msg, _ := mgr.Get("Norwegian Bokmål", err); msg != "Denne email eksisterer allerede i users." {
		t.Fail()
	}
}

func TestSimpleMultiple(t *testing.T) {
	mgr := NewManager([]string{"English", "Norwegian Bokmål"})

	errMsg := mgr.Add(`dial tcp: lookup port=: no such host`)
	errMsg.Set("English", "Unable to connect to external service.")
	errMsg.Set("Norwegian Bokmål", "Kunne ikke koble til ekstern tjener.")

	errMsg = mgr.Add(`pq: duplicate key value violates unique constraint "{{ table }}_{{ column }}_key"`)
	errMsg.Set("English", `This {{ column }} already exists in {{ table }}.`)
	errMsg.Set("Norwegian Bokmål", `Denne {{ column }} eksisterer allerede i {{ table }}.`)

	err := errors.New(`dial tcp: lookup port=: no such host`)
	if msg, _ := mgr.Get("English", err); msg != "Unable to connect to external service." {
		t.Fail()
	}

	err = errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)
	if msg, _ := mgr.Get("Norwegian Bokmål", err); msg != "Denne email eksisterer allerede i users." {
		t.Fail()
	}
}

var translations = map[string]map[string]string{
	"Norwegian Bokmål": map[string]string{
		"email": "eposten",
		"users": "brukere",
	},
}

var translateVariables = func(language string, data map[string]string) {
	for k, v := range data {
		transLang, ok := translations[language]
		if translation, ok2 := transLang[v]; ok && ok2 {
			data[k] = translation
			fmt.Println("Translated", translation)
		}
	}
}

func TestCustom(t *testing.T) {
	mgr := NewManager([]string{"English", "Norwegian Bokmål"})

	errMsg := mgr.Add(`pq: duplicate key value violates unique constraint "{{ table }}_{{ column }}_key"`)
	errMsg.Set("English", `This {{ column }} already exists in {{ table }}.`)
	errMsg.SetCustom("Norwegian Bokmål", `Denne {{ column }} eksisterer allerede i {{ table }}.`, translateVariables)

	err := errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)
	if msg, _ := mgr.Get("English", err); msg != "This email already exists in users." {
		t.Fail()
	}

	if msg, _ := mgr.Get("Norwegian Bokmål", err); msg != "Denne eposten eksisterer allerede i brukere." {
		t.Fail()
	}
}
