# sendx-assignment

## How to Run
1. install go
2. cd to main.go
3. go run main.go
4. Now open browser or postman and call localhost:portnumber/route?url=url&rate_limit


For example "http://localhost:8000/pagesource?url=https://www.facebook.com&retry_limit=3"
As you can see in this screenshot it took 1020ms to get response, because it's get the page and storing it


If you execute same request, it took 3ms because we already got the page and we are returning the data from cache

