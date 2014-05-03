hutch
=====

Get a daily report of trending links on Twitter on the topic(s) of your choice.


### Install


	go get -d github.com/njern/hutch
	cd $GOPATH/src/github.com/njern/hutch
	go build

### Usage

First, you'll need to fill out the `credentials.json` file according to the example in `credentials_example.json`. This will also require you to register your own Twitter "App" at [Twitter's developer page](https://dev.twitter.com/).

Then, simply start Hutch with a command like:

	./hutch --topics "Justin Bieber, Canada" --num_links 10

In twenty-four hours (and every twenty-four hours after that) you will receive an e-mail with the top 10 most talked about links matching one or more of your defined keywords.

### TODO:

* Add notification options other than e-mail
* Write some unit tests

### Contributors
* [njern](https://github.com/njern)
* [mkaz](https://github.com/mkaz)
