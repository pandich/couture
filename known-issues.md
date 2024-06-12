KNOWN ISSUES
============

WINDOWS
-------

### Terminal Window Auto-Resizing

Terminal resize events currently are handled by trapping the `SIGWINCH` signal. This signal is not available on Windows.
This functionality will be added to Windows in a future release.
