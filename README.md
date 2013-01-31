#watcher

Watches files and runs an input command on them when the file is changed

##Usage
```bash
watcher -c "cat" app.js
```

Run watcher on the command line, specify the command with the -c flag and then a file or list of files.

When the file is modified, the command will be run only on the file modified, leaving all others in their current condition.
