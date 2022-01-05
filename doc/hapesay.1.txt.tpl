hapesay(1)
=========

Name
----
Neo hapesay/hapethink - configurable speaking/thinking hape (and a bit more)

SYNOPSIS
--------
hapesay [-e _eye_string_] [-f _hapefile_] [-h] [-l] [-n] [-T _tongue_string_] [-W _column_] [-bdgpstwy]
       [--random] [--bold] [--rainbow] [--aurora] [--super] [_message_]

DESCRIPTION
-----------
_Neo-hapesay_ (hapesay) generates an ASCII picture of a hape saying something provided by the
user.  If run with no arguments, it accepts standard input, word-wraps
the message given at about 40 columns, and prints the hape saying the
given message on standard output.

To aid in the use of arbitrary messages with arbitrary whitespace,
use the *-n* option.  If it is specified, the given message will not be
word-wrapped.  This is possibly useful if you want to make the hape
think or speak in figlet(6).  If *-n* is specified, there must not be any command-line arguments left
after all the switches have been processed.

The *-W* specifies roughly (where the message should be wrapped. The default
is equivalent to *-W 40* i.e. wrap words at or before the 40th column.

If any command-line arguments are left over after all switches have
been processed, they become the hape's message. The program will not
accept standard input for a message in this case.

There are several provided modes which change the appearance of the
hape depending on its particular emotional/physical state. 

The *-b* option initiates Borg mode

*-d* causes the hape to appear dead 

*-g* invokes greedy mode

*-p* causes a state of paranoia to come over the hape

*-s* makes the hape appear thoroughly stoned

*-t* yields a tired hape

*-w* is somewhat the opposite of *-t* and initiates wired mode

*-y* brings on the hape's youthful appearance.

The user may specify the *-e* option to select the appearance of the hape's eyes, in which case
the first two characters of the argument string _eye_string_ will be used. The default eyes are 'oo'. The tongue is similarly
configurable through *-T* and _tongue_string_; it must be two characters and does not appear by default. However,
it does appear in the 'dead' and 'stoned' modes. Any configuration
done by *-e* and *-T* will be lost if one of the provided modes is used.

The *-f* option specifies a particular hape picture file (``hapefile'') to
use. If the hapefile spec contains '/' then it will be interpreted
as a path relative to the current directory. Otherwise, hapesay
will search the path specified in the *COWPATH* environment variable. If *-f -* is specified, provides
interactive Unix filter (command-line fuzzy finder) to search the hapefile.

To list all hapefiles on the current *COWPATH*, invoke *hapesay* with the *-l* switch.

*--random* pick randomly from available hapefiles

*--bold* outputs as bold text

*--rainbow* and *--aurora* filters with colors an ASCII picture of a hape saying something

*--super* ...enjoy!

If the program is invoked as *hapethink* then the hape will think its message instead of saying it.

COWFILE FORMAT
--------------
A hapefile is made up of a simple block of *perl(1)* code, which assigns a picture of a hape to the variable *$the_hape*.
Should you wish to customize the eyes or the tongue of the hape,
then the variables *$eyes* and *$tongue* may be used. The trail leading up to the hape's message balloon is
composed of the character(s) in the *$thoughts* variable. Any backslashes must be reduplicated to prevent interpolation.
The name of a hapefile should end with *.hape ,* otherwise it is assumed not to be a hapefile. Also, at-signs (``@'')
must be backslashed because that is what Perl 5 expects.

ENVIRONMENT
-----------
The COWPATH environment variable, if present, will be used to search
for hapefiles.  It contains a colon-separated list of directories,
much like *PATH or MANPATH*. It should always contain the */usr/local/share/hapes*
directory, or at least a directory with a file called *default.hape* in it.

FILES
-----
*%PREFIX%/share/hapes* holds a sample set of hapefiles. If your *COWPATH* is not explicitly set, it automatically contains this directory.

BUGS
----
https://github.com/Rid/hapesay

If there are any, please report bugs and feature requests in the issue tracker.
Please do your best to provide a reproducible test case for bugs. This should
include the *hapesay* command, the actual output and the expected output.

AUTHORS
-------
Neo-hapesay author is Kei Kamikawa (x00.x7f.x86@gmail.com).

The original author is Tony Monroe (tony@nog.net), with suggestions from Shannon
Appel (appel@CSUA.Berkeley.EDU) and contributions from Anthony Polito
(aspolito@CSUA.Berkeley.EDU).

SEE ALSO
--------
perl(1), wall(1), nwrite(1), figlet(6)