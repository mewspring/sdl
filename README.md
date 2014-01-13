WIP
---

This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation can be inaccurate.

win
===

Package win provides the core functionality required for window creation,
drawing and event handling. The window events are defined in a dedicated package
located at:
	github.com/mewmew/we

The library uses a small subset of the features provided by [SDL][libsdl]
version 2.0. Support for multiple windows has intentionally been left out to
simplify the API.

[libsdl]: http://www.libsdl.org/

Documentation
-------------

Documentation provided by GoDoc.

- sdl
   - [win][sdl/win]: handles window creation, drawing and events.

[sdl/win]: http://godoc.org/github.com/mewmew/sdl/win

Installation
------------

Install the [SDL][libsdl] library version 2.0 and run:

	go get github.com/mewmew/sdl/win

Examples
--------

simple demonstrates how to draw surfaces using the Draw and DrawRect methods. It
also gives an example of a basic event loop.

	go get github.com/mewmew/sdl/examples/simple

![Screenshot - simple](https://raw.github.com/mewmew/sdl/master/examples/simple/simple.png)

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
