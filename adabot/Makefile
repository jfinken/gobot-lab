all:
	GOARM=7 GOARCH=arm GOOS=linux go build -v adabot.go
	scp adabot pi@192.168.0.2:~
