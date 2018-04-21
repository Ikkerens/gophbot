# GophBot [![Build Status](https://travis-ci.org/ikkerens/gophbot.svg?branch=master)](https://travis-ci.org/ikkerens/gophbot) [![Go Report Card](https://goreportcard.com/badge/github.com/ikkerens/gophbot)](https://goreportcard.com/report/github.com/ikkerens/gophbot) 

This is a powerful moderation bot for use in Discord.
You can invite me using this link: [Invite me!](https://discordapp.com/oauth2/authorize?client_id=436984629027667968&scope=bot&permissions=0)

## Building
*All steps provided here assume you have your GOPATH configured correctly, see [this page](https://github.com/golang/go/wiki/SettingGOPATH).*

The two commands below will download all dependencies and then install the gophbot binary to `$GOPATH/bin`
```
go get github.com/ikkerens/gophbot/cmd/gophbot
go install github.com/ikkerens/gophbot/cmd/gophbot
```

## Running
Windows:
```
SET TOKEN="PASTEYOURTOKENHERE"
%GOPATH%\bin\gophbot.exe
```

Linux/Mac/Unix:
```
export TOKEN=PASTEYOURTOKENHERE
$GOPATH/bin/gophbot
```