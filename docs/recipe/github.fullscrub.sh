#!/bin/sh

git checkout --orphan newBranch
git add -A  # Add all files and commit them
git commit
git branch -D master  # Deletes the master branch
git branch -m master  # Rename the current branch to master

git remote remove origin
git remote add origin git@github.com:whereiskurt/tiogo.git

git push -f origin master  # Force push master branch to github
git gc --aggressive --prune=all     # remove the old files

git branch --set-upstream-to=origin/master master

##export PS1="\[\e]0;whereiskurt@\h: \w\a\]${debian_chroot:+($debian_chroot)}\[\033[01;32m\]whereiskurt@gopherit\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$"