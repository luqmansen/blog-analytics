# blog-analytics

Simple Go App with hexagonal architecture for web analytics

### Why
My domain currently expired, so I didn't get nice metrics from Cloudflare about traffic for my Blog anymore

### What recorded data
- Visited URL
- IP address
- Visited time
- and few more (json example below)

### What metrics could I get
- Daily, weekly, yearly visitor counter
- Most city / country visit
- Basically anything

### Dev Requirement
- go 1.14
- MongoDB (or just use docker compose)
- Docker
- docker-compose


### Endpoints
**POST** ``api/analytics`` <br>
JSON body post example

```json 
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

### Client-side example
````javascript

// Change to your backend
const URL = "http://localhost:8080"

function httpGetAsync(theUrl, callback)
{
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(xmlHttp.responseText);
    }
    xmlHttp.open("GET", theUrl, true); 
    xmlHttp.send(null);
}

postData = async (data) => {
    const payload = {
        "url" :window.location.href,
         "info" : JSON.parse(data)
    }
    const settings = {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
        }
    };
    try {
        const fetchResponse = await fetch(`${URL}/api/analytics`, settings);
        const data = await fetchResponse.json();
        return data;
    } catch (e) {
        return e;
    }    

}

(function() {
    httpGetAsync("https://ipinfo.io/json", postData)
})();
````
