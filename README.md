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

## customization

Instead of looking for a `scripts` dir, you can set an environment variable `SCRIPTY_DIR` to the name
of a folder, without slashes, to look for scripts in.

## shell completions

Execute one of the following snippets to add completion for your favorite shell. If you downloaded the the binary only, you'll want to use the following technique to download the completion script instead of `cp`ing it.

 ```shell
curl https://raw.githubusercontent.com/steventlamb/scripty/master/completions/scripty.sh > ~/.bash_completion.d/scripty
```

### bash

```shell
cp completions/scripty.sh ~/.bash_completion.d/scripty
```

### fish

```shell
cp completions/scripty.fish ~/.config/fish/completions/scripty.fish
```

## motivation

Scripty is inspired by the behavior of modern tools like fabric, vagrant, etc.
These tools allow you to run commands from any subdirectory of a project, which is very
convenient. I had gotten used to using fabric, and after switching to a folder full of
scripts for simplicity, I missed this functionality.

