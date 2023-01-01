# (GO-Talk) Chat Application
Implemented in Go
~~~
GoTalk Server
/////////////.......................................................................................
Commands:
=====> run                                       : Run GoTalk server
         *Options:
           -b [ip]                               : Bind server to given ip address, default=0.0.0.0
           -p [port]                             : Specify the connection port, default=8080
           -db [database-name]                   : Specify database name, default=go-talk
           -d                                    : enable debug mode
               *Options:
                   1.leave empity will build tests
                   2."l"=(lazy debug run) mean's only run server without building tests
                   3."lazy"=equal to option (1)
=====> migrate                                   : Migrate all models to database
         *Options:
           -db [database-name]                   : Specify database name, default=go-talk
=====> --help , -h, help                         : Print available commands
=====> -v , --version, version                   : Getting current version of GoTalk
////////////////////////////////////////////////////////////////////////////////////////////////////
~~~

# Test
Run server on debug mode to build and get test client in browser.
~~~
go run main.go run -d
~~~
test client would be accessible on [http://localhost:8080/test/](http://localhost:8080/test/)


* Note: Make sure your reverse proxy doesn't hide the user's real ip address.