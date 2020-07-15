## forego
+ https://github.com/ddollar/forego/pull/105
+ https://github.com/ddollar/forego/pull/124

<a href="https://circleci.com/gh/bugficks/forego">
  <img align="right" src="https://circleci.com/gh/bugficks/forego.svg?style=svg">
</a>

[Foreman](https://github.com/ddollar/foreman) in Go.

### Installation

[Downloads](https://github.com/bugficks/forego/releases)

##### Compile from Source

    $ go get -u github.com/bugficks/forego

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO

Use `forego help` to get a list of available commands, and `forego help
<command>` for more detailed help on a specific command.

### License

Apache 2.0 &copy; 2015 David Dollar
