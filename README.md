# ghost-project-finder
Scanner to find "rogue" project files on local workstations

# How to build it

You'll need to download the Go tools from https://golang.org/doc/install, or else use the Docker image at https://hub.docker.com/_/golang?tab=tags&page=1&ordering=last_updated.

With this installed, simply: 

```bash
$ go build
```

This will create a `ghost-project-finder` binary which will be suitable for whatever you compiled it on (so if you ran `go build` on a Mac, then you'll have a binary
suitable for a Mac.  Just copy and run, no libraries required!

# How to run it

```
$ ./ghost-project-finder -start /System/Volumes/Data -nosend
```

This will scan from the root path (well, the Mac mount where the writable root goes and the root folders symlink to) and will show you on-screen
any Premiere, Prelude, AfterEffects or Cubase files it finds.
