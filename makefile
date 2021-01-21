dev:
	go run *.go
table:
	sqlite3 message.sqlite3 < base.sql

