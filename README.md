# iRCornerNames

First project in Go.
Lots of things to improve.
Lots of things to add.
But it works for now.

It boots up a webserver that you can then connect to via http://localhost:8080 (can be changed with CLI argument `-port 8081`).
This is useful if you want to include that into your OBS Scene and let the viewer show which corner you are currently in.

Additional CLI arguments are:
- `-port 8081` to use a different port, if needed. Default is 8080
- `-debug` to show debugging messages
- `-webdebug` to show debugging messages related to the webserver
- `-showlapdist` to show additionally the distance traveled from S/F Line
