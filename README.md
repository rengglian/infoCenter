# infoCenter
Simple GO program to get some information from a DIR-200 internt radio and use it in kiosk mode on a Rapsberry Pi 3.
Optimized for RASPBERRY Pi Touch Display, 7" 800x480 screen resolution.

Page is also accessible via web server within network on port :8080

Layout inspired by LCARS (Font from https://gtjlcars.de)

Crosscompile for Raspberry Pi 3: `GOOS=linux GOARCH=arm GOARM=7 go build infoCenter.go` 

![alt text](https://raw.githubusercontent.com/rengglian/infoCenter/master/screen.png)

