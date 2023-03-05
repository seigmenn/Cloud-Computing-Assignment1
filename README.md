# assignment-1



## Endpoints

The service provides three different handlers that can be found at these endpoints:
```
/unisearcher/v1/uniinfo/
/unisearcher/v1/neighbourunis/
/unisearcher/v1/diag/
```



## Retrieve information for a given university
This endpoint retrieves information about a university given its name. Note that the name of the university may be partial or complete, and that several universities may match a search.

```
Path: uniinfo/{:partial_or_complete_university_name}/
```

The output of this endpoint will look something like this:

```
[
    {
        "name": "Harvard University",
        "country": "United States",
        "alpha_two_code": "US",
        "web_pages": [
            "http://www.harvard.edu/"
        ],
        "languages": {
            "eng": "English"
        },
        "maps": {
            "openStreetMaps": "https://www.openstreetmap.org/relation/6430384"
        }
    }
]OK

```

## Diagnostics interface
The diagnostic interface returns the statuscode from the two APIs, which version of the handler it is and how long the webservice has been active.

```
Path: unisearcher/v1/diag/
```

The diagnostics interface will look something like this:

```
{
  "universitiesapi": 200,
  "countriesapi": 200,
  "version": "v1",
  "uptime": 79.1989281
}OK
```

## Retrieving universities from neighbouring countries with a similar name

This endpoint retrieves information about universities from neighbouring countries that share a similar name. Note that the name may be partial or complete. In addition, one can set a limit for how many universities are retrieved.

```
Path: unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}
```

The output of this endpoint will look something like this:

```
[
        {
                "name": "Häme University of Applied Sciences",
                "country": "Finland",
                "alpha_two_code": "FI",
                "web_pages": [
                        "https://www.hamk.fi/"
                ],
                "languages": {
                        "fin": "Finnish",
                        "swe": "Swedish"
                },
                "maps": {
                        "openStreetMaps": "openstreetmap.org/relation/54224"
                }
        },
        {
                "name": "Central Ostrobothnia University of Applied Sciences",
                "country": "Finland",
                "alpha_two_code": "FI",
                "web_pages": [
                        "http://www.cou.fi/"
                ],
                "languages": {
                        "fin": "Finnish",
                        "swe": "Swedish"
                },
                "maps": {
                        "openStreetMaps": "openstreetmap.org/relation/54224"
                }
        }
]OK
```

