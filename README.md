Progress is a command-line task management tool.

Progress stores your task history in a SQL database and provides an interface
for adding, removing, and listing information about tasks which are currently
being worked on. Progress uses SQLite by default, but you could provide a
different database if you felt it was necessary to do so.

This is pretty roughly implemented right now, but it does work.

## Usage

Here are some example usage commands:

    $ prg @project add Metanic MTC https://metanic.org
    Adding project: Metanic (MTC) - https://metanic.org
    
    $ prg @project list
    [MTC] Metanic - https://metanic.org
    
    $ prg Testing out prg commands
    Adding task to Metanic: Testing out prg commands
    
    $ prg @project add "Another Project" _ Some other project. Doesn\'t matter, tbh.
    Adding project: Another Project (AOT) - Some other project. Doesn't matter, tbh.
    
    $ prg Not much.
    Adding task to Another Project: Not much.
    
    $ prg '~MTC' Fixing CORS issue
    Adding task to Metanic: ~MTC Fixing CORS issue
    
    $ prg @task list
    Another Project:
    - Not much. [2]
    Metanic:
    - ~MTC Fixing CORS issue [3]
    - Testing out prg commands [1]
    
    $ prg '~MTC' Breaking CORS cuz YOLO
    Adding task to Metanic: Breaking CORS cuz YOLO
    
    $ prg meow
    Adding task to Another Project: meow
    
    $ prg '~MTC' 'Eating M&Ms'
    Adding task to Metanic: Eating M&Ms
    
    $ prg @project list
    [MTC] Metanic - https://metanic.org
    [AOT] Another Project - Some other project. Doesn't matter, tbh.

## Roadmap

* [ ]: Basic tools for managing projects
  * [x]: Add projects
  * [x]: List projects
  * [ ]: Remove projects (should we even do this?)
  * [ ]: Ability to provide description (optionally, using $EDITOR)
* [ ]: Basic tools for managing tasks
  * [x]: Add tasks
  * [x]: List tasks
  * [ ]: Ability to provide long-form details (optionally, using $EDITOR)
  * [ ]: Remove tasks (should we even do this?)
* [ ]: Tagging tasks, projects, etc with arbitrary metadata
* [ ]: Support GPG Signing of database entries with every update
