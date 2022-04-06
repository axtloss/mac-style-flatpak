# mac-style-flatpak
A little daemon that watches over a specific directory (Probably an Applications directory) and installs any flatpakref that get's put inside it. <br>
this way you can do app installation similiar to how it's done on macos, just drag the app into the Applications folder and it's installed, if you want to remove it just move it out of the Applications folder(not implemented yet).

# Building
to build it you have to do 
```
git clone https://github.com/axtloss/mac-style-flatpak
cd mac-style-flatpak
go build
```
I wouldn't recommend installing that yet (`go install`) as it's not finished.

# Usage
Once again, this isn't finished, you shouldn't rely on it working 100%.
first you have to make a directory where you want to put the applications in to, I'll take `~/Applications` as an example
then you want to set the envvar `APPLICATIONS_PATH="$HOME/Applications/"`
then just execute the built binary.
Here's a one-liner:
`mkdir ~/Applications && APPLICATIONS_PATH="$HOME/Applications/ ./m`

# Supported file managers
yes, I sadly have to say that some file managers won't fully work (thanks gnome) <br>
basically, every file manager will support creating the directory and installing/deleting the files, but nautilus, the gnome file manager, won't allow you to run the generated .desktop file from the directory, so just double-clicking on the installed file won't execute it. <br>
Every other file manager should work with no issues.

# Roadmap
- [x] Support installing flatpak files <br>
- [] Support removing installed appr <br>
- [] alert user when application finished installing/removing <br>
- [] Support custom Applications directories <br>
