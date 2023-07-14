EmuTOS - 192 KB versions

These ROMs are suitable for the following hardware:
- ST / STf
- Mega ST
- emulators of the above

Unlike other ROM versions, they do NOT autodetect extra hardware, and might not
work on machines with additional hardware. For example, they don't work under
the Hatari emulator's Falcon emulation due to missing VIDEL support.
Also, they only work with a plain 68000 CPU.

Each ROM contains a single language:

etos192cz.img - Czech (PAL)
etos192de.img - German (PAL)
etos192es.img - Spanish (PAL)
etos192fi.img - Finnish (PAL)
etos192fr.img - French (PAL)
etos192gr.img - Greek (PAL)
etos192hu.img - Hungarian (PAL)
etos192it.img - Italian (PAL)
etos192nl.img - Dutch (PAL)
etos192no.img - Norwegian (PAL)
etos192pl.img - Polish (PAL)
etos192ru.img - Russian (PAL)
etos192se.img - Swedish (PAL)
etos192sg.img - Swiss German (PAL)
etos192tr.img - Turkish (PAL)
etos192us.img - English (NTSC)
etos192uk.img - English (PAL)

The following optional files are also supplied:
emuicon.rsc - contains additional icons for the desktop
emuicon.def - definition file for the above

Note that the emuicon.rsc file format differs from deskicon.rsc used by later
versions of the Atari TOS desktop.

Note that selecting Norwegian/Swedish currently sets the language to English,
but the keyboard layout to Norwegian/Swedish.

Due to size limitations, the 192 KB ROMs contain:
- no EmuCON
- limited desktop features (comparable to Atari TOS 1)
- no builtin text file viewer and print function

These ROM images have been built using:
make all192

This release has been built on Linux Mint (a Ubuntu derivative), using
Vincent Rivière's GCC 4.6.4 cross-compiler.  The custom tools used in
the build process were built with native GCC 4.8.4.

The source package and other binary packages are available at:
https://sourceforge.net/projects/emutos/files/emutos/1.2.1/

An online manual is available at the following URL:
https://emutos.github.io/manual/

The extras directory (if provided) contains:
(1) one or more alternate desktop icon sets, which you can use to replace
    the builtin ones.  You can use a standard resource editor to see what
    the replacement icons look like.
    To use a replacement set, move or rename the existing emuicon.rsc &
    emuicon.def files in the root directory, then copy the files containing
    the desired icons to the root, and rename them to emuicon.rsc/emuicon.def.
(2) a sample mouse cursor set in a resource (emucurs.rsc/emucurs.def).  This
    set is the same as the builtin ones, but you can use it as a basis to
    create your own mouse cursors.
    To use a replacement set, copy the files containing the desired mouse
    cursors to the root, and rename them to emucurs.rsc/emucurs.def.
For further information on the above, see doc/emudesk.txt.

If you want to read more about EmuTOS, please take a look at these files:

doc/announce.txt      - Introduction and general description, including
                        a summary of changes since the previous version
doc/authors.txt       - A list of the authors of EmuTOS
doc/bugs.txt          - Currently known bugs
doc/changelog.txt     - A summarised list of changes after release 0.9.4
doc/emudesk.txt       - A brief guide to the newer features of the desktop
doc/incompatible.txt  - Programs incompatible with EmuTOS due to program bugs
doc/license.txt       - The FSF General Public License for EmuTOS
doc/status.txt        - What is implemented and running (or not yet)
doc/todo.txt          - What should be done in future versions
doc/xhdi.txt          - Current XHDI implementation status

Additional information for developers (just in the source archive):

doc/install.txt       - How to build EmuTOS from sources
doc/coding.txt        - EmuTOS coding standards (never used :-) )
doc/country.txt       - An overview of i18n issues in EmuTOS
doc/fat16.txt         - Notes on the FAT16 filesystem in EmuTOS
doc/memdetect.txt     - Memory bank detection during EmuTOS startup
doc/nls.txt           - How to add a native language or use one
doc/old_changelog.txt - A summarised list of changes up to & including
                        release 0.9.4
doc/osmemory.txt      - All about OS internal memory in EmuTOS
doc/reschange.txt     - How resolution change works in the desktop
doc/resource.txt      - Modifying resources in EmuTOS
doc/startup.txt       - Some notes on the EmuTOS startup sequence
doc/tos14fix.txt      - Lists bugs fixed by TOS 1.04 & their status in EmuTOS
doc/version.txt       - Determining the version of EmuTOS at run-time

The following documents are principally of historical interest only:

doc/old_code.txt      - A museum of bugs due to old C language
doc/vdibind.txt       - Old information on VDI bindings

-- 
The EmuTOS development team
https://emutos.sourceforge.io/
