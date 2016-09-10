all:
	# Pi 3
	GOARM=7 GOARCH=arm GOOS=linux go build -v 
	# Pi 2
	#GOARM=6 GOARCH=arm GOOS=linux go build -v 
	scp robot pi@192.168.0.2:~/robot-svc
	scp -r html pi@192.168.0.2:~/
