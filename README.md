# scripty

Run scripts from a nested directory without navigating up to project root.

## usage

Say you have the follow directory tree for a project you are developing, 

```
~/project/
├── scripts
│   ├── bar.sh
│   ├── baz.py
│   ├── foo.sh
│   └── qux
└── src
    └── subpackage
```

where the contents of scripts are all executable.

You can execute these scripts, with or without file extensions, from any directory under `~/project/`:
```
[~/project/src/subpackage]$ scripty -l
bar
baz
foo
qux
[~/project/src/subpackage]$ scripty bar
# ... executes bar.sh
[~/project/src/subpackage]$ scripty bar.sh
# ... executes bar.sh

```

## motivation

Scripty is inspired by the behavior of modern tools like fabric, vagrant, etc.
These tools allow you to run commands from any subdirectory of a project, which is very 
convenient. I had gotten used to using fabric, and after switching to a folder full of
scripts for simplicity, I missed this functionality.

