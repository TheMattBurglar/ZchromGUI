# ZchromGUI
This is a GUI version of my Zchrom CLI app.  It is a program to test what would happen if we added a Z chromosome that made a person Female, regardless of their other chromosome. ex: a ZY Female

While the CLI app is complete, I always wanted to see if I could make a GUI version.  Now, using AI assistance, I've made good progress on that.

The CLI version is 100% AI free.  The GUI version is significantly AI assisted, at least for the GUI elements.  That’s why I’m treating them like different apps all together.  If you want to use non-AI code, stick with the CLI.

This has been updated to use fyne, which is a GUI framework for Go.  I've also added a web version, which is a bit more complex, but should be even easier to run.

Known issue:
This does not have the option of starting with a random population, but I never found that option produced interesting results, so I don’t really care enough to put it back in.


This time, running this code locally is a bit more complicated than the CLI, because of the fyne code I used for the GUI.  You’ll need a number of development libraries to get fyne working.  Or, if you're running a recent version of Linux, you can probably just download the main_gui binary file and run it.

To run the GUI from source, use:
`go run .`
or
`go run main_gui.go charts.go`


The web version, though, should work on any system with a web browser. Just run `go run web/main_web.go` and open http://localhost:8080 in your browser.

## License

- **Code:** MIT License
- **Short story:** Creative Commons Attribution 4.0 International (CC BY 4.0)
