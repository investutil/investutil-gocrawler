# web-scraping
The service of Investutil needs data. Data from web scraping of publicly accessible sites is a very important source, besides the data gathered from free or paid APIs.

## 
Architecture
We use Crawlab as web crawler management platform
https://github.com/crawlab-team/crawlab
Crawlab is a golang-based distributed web crawler management platform. It support crawler frameworks including Scrapy, Selenium.

The web scraping framwork, we will use Scrapy(Python) and Colly (Golang)
Why not rust:
Python and Go are easier to code, and web scraping requires more agility.
