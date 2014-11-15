#MONOLITH

The Monolith is a single Go app (deployed easily to Heroku or AWS) that provides useful utility services to your web applications. Monolith can:

###Fetch Pages [Beta]
Tell the Monolith to fetch a page and send you the content when it's done. Useful for crawling, scraping, feeds and load testing.

###Resize Images Dynamically [In Development]
Store your images on a fileserver or S3 and use Monolith URLs to access them in the image size and pixel density that you need. Images are resized on the fly. For best results, use a CDN to cache and distribute the resized images.

###Broadcast Messages [In Development]
The Monolith lets web browsers hold an open connection to it via SSE (Server Sent Events - similar to, but much more lightweight than WebSockets). Your app can then publish events to the topics, which are delivered to all listening subscribers. Great for sending your web users any kinds of notifications instead of having them poll your server, or setting up a complicated and less scalable WebSocket server cluster.

###Act as a Task Queue [In Development]
Just tell Monolith to call a URL and it will do so until it receives a successful response. Also schedule calls in the short term future to help with task management or rate limiting. Because of the ability to have the tasks call URLs on multiple hosts, it becomes trivial to have the Monolith manage a work pipeline or a distributed system as well.

##Reliability
Monolith is designed to be massively distributed, so all operations work on a best effort basis. In a multi-server setup, the Monolith requires one instance of Redis per cluster to co-ordinate work and record measurements. As always, it makes sense to design for failure, automatic retries and [idempotence](https://en.wikipedia.org/wiki/Idempotence). 
