# blog-analytics

Simple Go App with hexagonal architecture for web analytics

### Why
My domain currently expired, so I don't nice metrics from Cloudflare about traffic on my Blog anymore

### What recorded data
- Visited URL
- IP address
- Visited time
- etc

### Dev Requirement
- go 1.14
- Docker
- docker-compose


### Endpoints
**GET** ``api/analytics`` <br>

Success response example
```json
[
  {
    "created_at": "2020-08-02T07:32:37.996Z",
    "url": "https://luqmansen.github.io/yeetV2",
    "ip": "127.0.0.1:53780",
    "info": {
        "ip": "69.69.69.69",
        "city": "San Fransisco",
        "region": "California",
        "country": "US",
        "loc": "-69.240,69.420",
        "org": "Google",
        "timezone": "US/UTC"
    }
  }
]
```

**POST** ``api/analytics`` <br>
JSON body post example
````json 
{
	"url": "https://luqmansen.github.io/yeetV3",
	"info": {
		"ip": "69.69.69.69",
		"city": "San Fransisco",
		"region": "California",
		"country": "US",
		"loc": "-69.240,69.420",
		"org": "Google",
		"timezone": "US/UTC",
	}
}
```