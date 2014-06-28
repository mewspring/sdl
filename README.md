WIP
---

This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation can be inaccurate.

sdl
====

This project implements window creation, event handling and image drawing using
[SDL][] version 2.0.

[SDL]: https://www.libsdl.org/

Documentation
-------------

Documentation provided by GoDoc.

- sdl
   - [font][sdl/font]: handles graphical text entries with customizable font
   size, style and color.
   - [texture][sdl/texture]: handles hardware accelerated image drawing.
   - [window][sdl/window]: handles window creation, drawing and events.

[sdl/font]: http://godoc.org/github.com/mewmew/sdl/font
[sdl/texture]: http://godoc.org/github.com/mewmew/sdl/texture
[sdl/window]: http://godoc.org/github.com/mewmew/sdl/window

Examples
--------

### tiny

The [tiny][examples/tiny] command demonstrates how to render images onto the
window using the [Draw][sdl/window#Window.Draw] and
[DrawRect][sdl/window#Window.DrawRect] methods. It also gives an example of a
basic event loop.

	go get github.com/mewmew/sdl/examples/tiny

![Screenshot - tiny](https://raw.github.com/mewmew/sdl/master/examples/tiny/tiny.png)

[examples/tiny]: https://github.com/mewmew/sdl/blob/master/examples/tiny/tiny.go#L37
[sdl/window#Window.Draw]: http://godoc.org/github.com/mewmew/sdl/window#Window.Draw
[sdl/window#Window.DrawRect]: http://godoc.org/github.com/mewmew/sdl/window#Window.DrawRect

### fonts

The [fonts][examples/fonts] command demonstrates how to render text using TTF
fonts.

	go get github.com/mewmew/sdl/examples/fonts

![Screenshot - fonts](https://raw.github.com/mewmew/sdl/master/examples/fonts/fonts.png)

[examples/fonts]: https://github.com/mewmew/sdl/blob/master/examples/fonts/fonts.go#L39

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
