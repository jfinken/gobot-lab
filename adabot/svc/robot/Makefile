all:
	GOARM=7 GOARCH=arm GOOS=linux go build -v 
	scp robot pi@192.168.0.2:~/robot-svc
