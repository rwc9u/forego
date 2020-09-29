## forego

[Foreman](https://github.com/ddollar/foreman) in Go.

<a href="https://circleci.com/gh/rwc9u/forego">
  <img align="right" src="https://circleci.com/gh/rwc9u/forego.svg?style=svg">
</a>

This fork includes changes and updates from
* https://github.com/bugficks/foreman

including the following PRs against the original forego

* https://github.com/ddollar/forego/pull/105
* https://github.com/ddollar/forego/pull/124

Additional changes include reverse proxy support and other minor refactors.

### Installation

[Downloads](https://github.com/rwc9u/forego/releases)

##### Compile from Source

    $ go get -u github.com/rwc9u/forego

##### Brew tap

    $ brew tap rwc9u/forego
    $ brew install rwc9u/forego/forego

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO

    $ forego start -x 9999 -c web=3
    forego   | starting web.1 on port 5000
    forego   | starting web.2 on port 5001
    forego   | starting web.3 on port 5002
    forego   | Starting reverse proxy on port 9999
    web.1    | Example app listening at http://localhost:5000
    web.2    | Example app listening at http://localhost:5001
    web.3    | Example app listening at http://localhost:5002
    worker.1 | listening to queue FOO

Use `forego help` to get a list of available commands, and `forego help
<command>` for more detailed help on a specific command.

### License

Apache 2.0 &copy; 2015 David Dollar
