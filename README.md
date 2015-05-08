# config
---

A package that helps load an application's configuration through a mixture
of flags and environment variables.

## Purpose

Many applications take command line parameters or check environment
variables to figure out where their database is or what `S3` bucket
they need to throw files into. Depending on the use case, you may
only want flags, only environment variables or a mixture of both. This
package aims to provide a simple and straightforward manner of
achieving that.

**This project is still under development. The API is not guaranteed to stay the same.**