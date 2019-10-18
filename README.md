oibot
===
A Go library for interacting with iRobot® Create® 2 robots according to the [iRobot® Open Interface (OI) specification](https://www.irobotweb.com/~/media/MainSite/PDFs/About/STEM/Create/iRobot_Roomba_600_Open_Interface_Spec.pdf).


This library was based on the [github.com/xa4a/go-roomba](https://github.com/xa4a/go-roomba) project which has not been updated in several years, seemingly abandoned, yet surprisingly still functional. A large portion of that library has been removed so that only essential functionality is implemented.


That project was "remotely inspired" by the `pyrobot` library by damonkohler@gmail.com (Damon Kohler).


I've removed the simulator and Go test harness capabilities of the previous project for the sake of simplicity, as it wasn't immediately apparent how to use them or how complete their test coverage actually was.


Serial support is implemented with [github.com/tarm/serial](https://github.com/tarm/serial).


Dependencies
---
- Serial streams - [github.com/tarm/serial](https://github.com/tarm/serial)


License
---
This software is licensed to *you* under the terms of the MIT license


Usage
===
- TBD!

