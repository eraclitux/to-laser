=======
to-laser
=======

``tl`` is a command line gcode post processor for `jscut <https://github.com/tbfleming/jscut>`_ that prapares gcode for a ``Grbl`` controlled laser cutter.

Examples
========

Just redirect gcode via pipe::

        cat gcode.nc | tl -power 80 > laser.gcode
