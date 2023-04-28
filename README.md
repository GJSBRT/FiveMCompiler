# FiveMCompiler
Make it harder to read and exploit FiveM scripts. A bunch of optimizations are needed like using go routines. 
Do not run this if you do not use a version control system like Git. You should run this over your entire resource folder,
as it renames events.

## Functionality
- Renames events
- Minifies code

## Running
```sh
./fivemcompiler --dir "./resources"
```
