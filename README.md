# tasks

![](./tasks.gif)

## Quick Start

```sh
go install github.com/MoXcz/tasks@latest
```

And then run the program:

```sh
tasks --help
```

## Usage

```sh
tasks add "Do the dishes"
tasks add "Continue customizing Neovim (btw)"
tasks list
```

Output:
```
Total tasks: 2
ID   Task                                Created             Done
1    Do the dishes                       a few seconds ago   false
2    Continue customizing Neovim (btw)   a few seconds ago   false
```

```sh
tasks complete 1
tasks list # "Do the dishes" no longer appears
```

Output:
```
Total tasks: 2
ID   Task                                Created        Done
2    Continue customizing Neovim (btw)   a minute ago   false
```

```sh
tasks list --all # Show all tasks, including completed ones (which includes "Do the dieshes"
```

Output:
```
Total tasks: 2
ID   Task                                Created        Done
1    Do the dishes                       a minute ago   true
2    Continue customizing Neovim (btw)   a minute ago   false
```

## Build

```sh
git clone github.com/MoXcz/tasks.git
cd tasks
go build .
./tasks
```

## TODO

- [ ] Add `clean` command to remove all completed tasks. The intention would be to centralize how *done tasks* are handled to avoid having to read the tasks file in order to find the next `id` by removing the necessity to manually clean the file. The current implementation for searching the next ID is to just read the last record.
- [ ] Add `password` command (with things like `add` or `use`) to manage passwords. For now it's just an idea, it probably won't replace any existing password managers.
- [ ] Add `export` command to export tasks to a file in a specific format (JSON or CSV).
- [ ] Add `import` command to import tasks from a file in a specific format (JSON or CSV).
- [ ] Maybe add `tags`
- [x] Format `list` command: https://github.com/mergestat/timediff and text/tabwriter
- [ ] Add support for JSON and SQLite
- [ ] Refactor how the task file is handled when rewriting it.

Inspired by:
1. https://github.com/dreamsofcode-io/goprojects

