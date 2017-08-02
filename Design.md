Design of the locking algorithm
===============================

Problem scope
-------------

The design space of the lockfile library is to manage concurrent invocations of
unrelated machine local processes manipulating data locally to machine.


What is needed from the file system to support locking
------------------------------------------------------
 * suport for hard links in the filesystem. And if missing, reporting the absence of that
   support via an error telling it so, while trying to create hard links.
 * ability to atomically set hard links. If two processes try to set a hard link from file
   X to hard link Y, only one of them can do it and the other one will notice
   that it failed to set the hard link with an error detected by the Go function 
   os.IsExist as such an error.
 * atomic rename of files. If two processes try to rename name a file
   to the new name X, only one of them can do it and the other one will notice
   that it failed to set the new name to X with an error detected by the Go function 
   os.IsExist as such an error.
 * the filesystem subtree, where the lock file resides on, needs to be mounted only to one machine
   at one time between TryLock and Unlock.
