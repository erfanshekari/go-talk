# (GO-Talk) Chat Application
Implemented in Go
~~~
GoTalk Server Commands

-> run                        :   Run GoTalk server
-> --help , -h, help          :   Print available commands
-> -v , --version, version    :   Getting current version of GoTalk
-------------------------------------------------------------------
Flags
-b [ip]                       :    Bind server to given ip address
-p [port]                     :    set port
-d options("", "lazy")        :    enable debug mode
-------------------------------------------------------------------
~~~

# Test
Run server on debug mode to build and get test client in browser.
~~~
go run main.go run -d
~~~
test client would be accessible on [http://localhost:8080/test/](http://localhost:8080/test/)
