# sendx-assignment

## How to Run
1. install go
2. cd to main.go
3. go run main.go
4. Now open browser or postman and call localhost:portnumber/route?url=url&rate_limit


For example "http://localhost:8000/pagesource?url=https://www.facebook.com&retry_limit=3"
As you can see in this screenshot it took 1020ms to get response, because it's get the page and storing it
![first_time](https://user-images.githubusercontent.com/54117043/209856523-5bb7fd62-e717-488a-afab-feada291e1d4.jpeg)


If you execute same request, it took 3ms because we already got the page and we are returning the data from cache
![next_time](https://user-images.githubusercontent.com/54117043/209856541-e4938cf9-57f0-4676-b4d8-757f9334782e.jpeg)
