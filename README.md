# ZchromGUI
This is a GUI version of my Zchrom CLI app.  It is a program to test what would happen if we added a Z chromosome that made a person Female, regardless of their other chromosome. ex: a ZY Female

While the CLI app is complete, I always wanted to see if I could make a GUI version.  Unfortunately, I never had the time to learn how.  But now, using copilot, I was able to put together a GUI in an afternoon.

The CLI version is 100% AI free.  The GUI version is significantly AI assisted, at least for the GUI elements.  That’s why I’m treating them like different apps all together.  If you want to use non-AI code, stick with the CLI.

Known issues:
This does not have the option of starting with a random population, but I never found that option produced interesting results, so I don’t really care enough to put it back in.

Right now, there also isn’t an option to have viable Y chromosome eggs.  It didn’t produce the results I wanted for the story I was writing, so I rarely used it.  However, it is a useful option to have if you want to engage with the hypothetical Z chromosome idea from multiple angles.  I might get around to putting that option back in at some point.

The example output in the GUI won’t show the population totals or percentages if they lost all the Z chromosomes by the end, men died off, or the population hit the population cap.  This is due to the way I originally put together the code, resulting in the final runthrough in those instances changing the population to a marker instead of an actual total.  For that reason, I still recommend running this in a terminal, so you can see the output for the whole run, rather than just an example in the GUI.  I’m not sure if this is an issue I’m going to fix.

This time, compiling and running this code is a bit more complicated than the CLI, because of the fyne code I used for the GUI.  You’ll be needing a number of development libraries to get this working, but it shouldn’t be too hard to figure out.
