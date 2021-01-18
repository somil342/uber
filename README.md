# uber
3 Apis

1.See past booking

->http://localhost:8080/pastbooking
//header-> username

2.Find Cabs near by provided location

->http://localhost:8080/cabsnearbyme?loc=1,2&distance=10

3.Create a Ride for source to destination

->http://localhost:8080/bookcab
//header-> username
//json->{"fromrow":0,"fromcol":1,"torow":0,"tocol":4}
