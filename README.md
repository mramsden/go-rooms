# go-rooms

Forum software that allows people to post and reply to
messages publically.

## Running

To run the server:

```
go run ./cmd/server
```

## Development

### Organisation

Where should I put code? There are a few folders to be
aware of when adding code to this project. `cmd` is where
external commands live. `cmd/server` is the main HTTP server
that is responsible for serving HTML to users. It takes care
of routing and template rendering the data returned from
the core data layer. `internal` is for common code that is
re-used by application code in `cmd`. Core business logic
can live here for example, `users` includes code related to
managing user accounts and `testutils` includes code to
support writing tests.

`migrations` are where all of the database schema changes
live.

### Database migrations

Database migrations are run using [`golang-migrate`](https://github.com/golang-migrate/migrate).
When adding migrations to the project you must include both
an up and down migration following the numbered format.
