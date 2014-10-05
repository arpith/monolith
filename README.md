#MONOLITH

The Monolith is a single Go app (deployed easily to Heroku or AWS) that provides a fairly common set of utilities to web applications. Operations include:

###Fetch
Tell the Monolith to fetch a page and send you the content when it's done. Useful for crawling, scraping and load testing.

###Images
Store your images on a fileserver or S3 and use Monolith URLs to access them in the file size and pixel density that you need for the display context. Images are resized on the fly. For best results, use a CDN to cache and distribute the resized images.

###Broadcasting
The Monolith lets web browsers hold an open connection to it via HTTP SSE (Server Sent Events - similar to, but much more lightweight than Websockets). Your app can then publish events to the topics, which are delivered to all listening subscribers. Great for sending users any kinds of notifications instead of having them poll your server.

###Task Queue
The Monolith also acts as a lightweight task queue - just tell it to call a URL (either on the same source server or a different one) and it will do so until it receives a successful response. Also schedule calls in the short term future to help with task management or rate limiting. Because of the ability to have the tasks call URLs on multiple hosts, it becomes trivial to have the Monolith manage a pipeline as well.
