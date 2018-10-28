EmuTOS - 192 KB versions

These ROMs are suitable for the following hardware:
- ST / STf
- Mega ST
- emulators of the above

Unlike other ROM versions, they do NOT autodetect extra hardware, and might not
work on machines with additional hardware. For example, they don't work under
the Hatari emulator's Falcon emulation due to missing VIDEL support.

Each ROM contains a single language:

etos192cz.img - Czech (PAL)
etos192de.img - German (PAL)
etos192es.img - Spanish (PAL)
etos192fi.img - Finnish (PAL)
etos192fr.img - French (PAL)
etos192gr.img - Greek (PAL)
etos192it.img - Italian (PAL)
etos192no.img - Norwegian (PAL)
etos192ru.img - Russian (PAL)
etos192se.img - Swedish (PAL)
etos192sg.img - Swiss German (PAL)
etos192us.img - English (NTSC)
etos192uk.img - English (PAL)

The following optional files are also supplied:
emucurs.rsc - modifiable mouse cursors for the AES/desktop
emucurs.def - definition file for the above
emuicon.rsc - contains additional icons for the desktop
emuicon.def - definition file for the above

Note that the emuicon.rsc file format differs from deskicon.rsc used by later
versions of the Atari TOS desktop.

Note that selecting Norwegian/Swedish currently sets the language to English,
but the keyboard layout to Norwegian/Swedish.

Due to size limitations, the 192 KB ROMs contains:
- no EmuCON

Since EmuTOS 0.9.1, the 192 KB ROMs contain the same desktop as other variants.

These ROM images have been built using:
make all192

This release has been built on Windows/Cygwin, using
Vincent Rivière's GCC 4.6.4 cross-compiler.  The custom tools used in
the build process were built with native GCC 6.4.0.

The source package and other binary packages are available at:
http://sourceforge.net/projects/emutos/files/emutos/0.9.9.1/

If you want to read more about EmuTOS, please take a look at these files:

doc/announce.txt      - Introduction and general description, including
                        a summary of changes since the previous version
doc/authors.txt       - A list of the authors of EmuTOS
doc/bugs.txt          - Currently known bugs
doc/changelog.txt     - A list of changes: detailed up to and including
                        version 0.9.4; summarised for subsequent versions
doc/emudesk.txt       - A brief guide to the newer features of the desktop
doc/license.txt       - The FSF General Public License for EmuTOS
doc/license_aros.txt  - The AROS Public License for certain Amiga source
                        code distributed as part of EmuTOS source
doc/status.txt        - What is implemented and running (or not yet)
doc/todo.txt          - What should be done in future versions
doc/xhdi.txt          - Current XHDI implementation status

Additional information for developers (just in the source archive):

doc/install.txt       - How to build EmuTOS from sources
doc/coding.txt        - EmuTOS coding standards (never used :-) )
doc/country.txt       - An overview of i18n issues in EmuTOS
doc/fat16.txt         - Notes on the FAT16 filesystem in EmuTOS
doc/incompatible.txt  - Programs incompatible with EmuTOS due to program bugs
doc/memdetect.txt     - Memory bank detection during EmuTOS startup
doc/nls.txt           - How to add a native language or use one
doc/osmemory.txt      - All about OS internal memory in EmuTOS
doc/reschange.txt     - How resolution change works in the desktop
doc/resource.txt      - Modifying resources in EmuTOS
doc/tos14fix.txt      - Lists bugs fixed by TOS 1.04 & their status in EmuTOS

The following documents are principally of historical interest only:

doc/bios.txt          - The Hitchhiker's Guide to the BIOS (Atari document)
doc/old_code.txt      - A museum of bugs due to old C language
doc/vdibind.txt       - Old information on VDI bindings

-- 
The EmuTOS development team
http://emutos.sourceforge.net/
