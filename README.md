Go Copy Zip file With Filter
============================

Unsurprisingly, this was written for Log4shell remediation, to remove
the Jndi class from log4j-core.jar's

It gives you just 1 function:

func CopyZipWithoutFile(origPath string, skipFileRE *regexp.Regexp, newSuffix string) error 

Given a path, it assumes the file is a jar file it copies the file to
a new file named the same plus newSuffix

Any paths that match skipFileRE will not be copied.

The sizes seem a bit off, perhaps different Deflate implementations?


Added a

func CompareZipFiles(sourcePath string, destPath string, skipFileRE *regexp.Regexp)

to validate that everything in sourcePath intended is in destPath with
same SHA and FileHeader