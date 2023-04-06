# cdf: Cubeflix Document Format

A CDF document is structured as a tree, containing *blocks*. A *block*, like a paragraph or image, may contain *inline blocks* or more blocks.

## Pages

CDF pages is a server that allows for the hosting and creation of CDF documents. A CDF pages project contains the following files:

* `template.html`: the base page template
* `404.html`: 404 template
* `invalid.html`: invalid page template
* `pages/`: all page sources
* `pages/index.cdf`: index page
* `static/`: static files

## Todo

* file editing/live update