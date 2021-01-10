# TODO

- [ ] Switch to using `filepath.Walk` for path-based scanning
- [ ] Try resolving UNKNOWN symbols by examining Imported Libraries
    - [ ] Read in exported symbols for each library into a Links structure
    - [ ] Try resolving against each new Links structure
- [ ] Make sure that abi-wizard can be easily imported and used as a library
    - [ ] Hyphens allowed in package name?
        - [ ] Sub-packages?
    - [ ] Missing public variables?
    - [ ] Anything missing from the interface?
