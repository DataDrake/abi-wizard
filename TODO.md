# TODO

- [x] Switch to using `filepath.Walk` for path-based scanning
- [x] Try resolving UNKNOWN symbols by examining Imported Libraries
    - [x] Read in exported symbols for each library into a Links structure
    - [x] Try resolving against each new Links structure
- [x] Make sure that abi-wizard can be easily imported and used as a library
    - [x] Hyphens allowed in package name?
        - [x] Sub-packages?
    - [x] Missing public variables?
    - [x] Anything missing from the interface?
