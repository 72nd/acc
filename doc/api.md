---
title: ACC REST API Documentation v
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="">ACC REST API Documentation v</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

This API enables to use the plain-text ERP tool Acc via an REST interface.

Base URLs:

* <a href="https://localhost/">https://localhost/</a>

Email: <a href="mailto:msg@frg72.com">72nd</a> Web: <a href="https://github.com/72nd">72nd</a> 
License: <a href="https://opensource.org/licenses/MIT">MIT</a>

<h1 id="-default">Default</h1>

## get__customers

> Code samples

```shell
# You can also use wget
curl -X GET https://localhost/customers

```

```http
GET https://localhost/customers HTTP/1.1
Host: localhost

```

```javascript

fetch('https://localhost/customers',
{
  method: 'GET'

})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

```ruby
require 'rest-client'
require 'json'

result = RestClient.get 'https://localhost/customers',
  params: {
  }

p JSON.parse(result)

```

```python
import requests

r = requests.get('https://localhost/customers')

print(r.json())

```

```php
<?php

require 'vendor/autoload.php';

$client = new \GuzzleHttp\Client();

// Define array of request body.
$request_body = array();

try {
    $response = $client->request('GET','https://localhost/customers', array(
        'headers' => $headers,
        'json' => $request_body,
       )
    );
    print_r($response->getBody()->getContents());
 }
 catch (\GuzzleHttp\Exception\BadResponseException $e) {
    // handle exception or api errors.
    print_r($e->getMessage());
 }

 // ...

```

```java
URL obj = new URL("https://localhost/customers");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "https://localhost/customers", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

`GET /customers`

*Returns all customers.*

<h3 id="get__customers-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|default|Default|Default response|None|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_Parties">Parties</h2>
<!-- backwards compatibility -->
<a id="schemaparties"></a>
<a id="schema_Parties"></a>
<a id="tocSparties"></a>
<a id="tocsparties"></a>

```json
[
  {
    "id": "string"
  }
]

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Party](#schemaparty)]|false|none|none|

<h2 id="tocS_Party">Party</h2>
<!-- backwards compatibility -->
<a id="schemaparty"></a>
<a id="schema_Party"></a>
<a id="tocSparty"></a>
<a id="tocsparty"></a>

```json
{
  "id": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|UUID of the party object|

