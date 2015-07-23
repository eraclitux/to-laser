=======
toLaser
=======

toLaser is a command line gcode post processor for ``jscut`` that prapares gcode for a ``Grbl`` controlled laser cutter.

Examples
========

Just redirect gcode via pipe::

        cat gcode.nc | toLaser > laser.gcode
