This is Break Out Of Hareka's Underground (short BOOHU), a roguelike game which
takes inspiration mainly from DCSS and its tavern, and some ideas from Brogue,
but aiming for very short games, almost no character building, and simplified
inventory management.

It is a work in progress, but is already a quite complete game.

Install
-------

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`).
+ Add `$GOPATH/bin` to your `$PATH`.
+ Use the command `go get github.com/anaseto/gofrundis/bin/frundis`.
  
The `boohu` command should now be available.

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.
