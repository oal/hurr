# hurr - Human Errors
`hurr` is a library for converting internal error messages to more user friendly versions.

An end user shouldn't have to encounter `pq: duplicate key value violates unique constraint "users_email_key"` in your application.

With `hurr` you can do this:

```go
	// Create an error manager
	mgr := hurr.NewManager([]string{"English", "Norwegian Bokmål"})

	// Create an error template
	errMsg := mgr.Add(`pq: duplicate key value violates unique constraint "{{ table }}_{{ column }}_key"`)

	// Add messages in whatever languages your application supports, reusing the variables extracted from the error template above
	errMsg.Set("English", `This {{ column }} already exists in {{ table }}.`)
	errMsg.Set("Norwegian Bokmål", `Denne {{ column }} eksisterer allerede i {{ table }}.`)

	// Somewhere in your app, this error may be returned from PostgreSQL
	err := errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)

	// Get the English end user friendly version of this error
	msg, _ := mgr.Get("English", err)
	log.Println(msg)
	// -> This email already exists in users.
```

Add as many error messages as you need to a `hurr.Manager`, and get their corresponding "human friendly" versions by calling `mgr.Get(language, error)`, where `language` is one of your manager's supported languages, and `error` is the `error` returned from some third party library. `hurr` will scan through the error and look for a match, and extracting any variables, between `{{` and `}}`.

See `hurr_test.go` for more examples, as well as how to use `.SetCustom`, and providing custom translations for variables extracted from error messages.