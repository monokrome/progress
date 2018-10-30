progress
--------

[![CircleCI](https://circleci.com/gh/monokrome/progress.svg?style=svg)](https://circleci.com/gh/monokrome/progress)

Progress is a command-line task management tool.

Progress stores your task history in a SQL database and provides an interface
for adding, removing, and listing information about tasks which are currently
being worked on. Progress uses SQLite by default, but you could provide a
different database if you felt it was necessary to do so.

This is pretty roughly implemented right now, but it does work.

## Usage

Here are some example usage commands:

    $ prg project create Metanic -a MTC
    Adding project: Metanic [MTC]

    $ prg project list
    [MTC] Metanic

    $ prg task create -a MTC Testing out prg commands
    Created task in Metanic: [MTC] Testing out prg commands

    $ prg project create Another Project
    Created project: Another Project [AOT]

    $ prg task create Not much.
    Created task in Another Project: Not much.

    $ prg create task -a MTC Fixing CORS issue
    Created task in Metanic: [MTC] Fixing CORS issue

    $ prg task tag fancy
    [MTC] Fixing CORS issue @fancy

    $ prg task tag done
    [MTC] Fixing CORS issue @fancy @done

    $ prg task tag -d fancy
    [MTC] Fixing CORS issue @done

    $ prg task list
    Another Project
    - Not much.

    Metanic
    - Fixing CORS issue
    - Testing out prg commands

    $ prg -a MTC Breaking CORS cuz YOLO
    Created task in Metanic: Breaking CORS cuz YOLO

    $ prg create task meow
    Created task in Another Project: meow

    $ prg @project list
    [MTC] Metanic
    [AOT] Another Project


## Roadmap

* [ ]: Basic tools for managing projects
  * [x]: Add projects
  * [x]: List projects
  * [x]: Remove projects
  * [ ]: Ability to provide description (optionally, using $EDITOR)
* [ ]: Basic tools for managing tasks
  * [x]: Add tasks
  * [x]: List tasks
  * [ ]: Ability to provide long-form details (optionally, using $EDITOR)
* [x]: Tagging tasks, projects, etc with arbitrary metadata
* [ ]: Support GPG Signing of database entries with every update
