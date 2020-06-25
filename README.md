# zendesk-holiday-creator
Bulk creates Zendesk holidays for schedules

## Use

Clone this repository and run `go build`. To run, ensure that `holidays.txt` is up to date with a comma separated list of all of the holidays to be added. In your command line, run:

```
./zendesk-holiday-updater -url="<ZENDESK-SUBDOMAIN>" -id="<SCHEDULE-ID>" -user="<ZENDESK-API-USER>" -key="<ZENDESK-API-KEY>"
```

You should then see the results of each holiday logged in your terminal. 

```
&{POST https://<SUBDOMAIN>.zendesk.com/api/v2/business_hours/schedules/401068/holidays.json HTTP/1.1 1 1 map[] {{"holiday": {"name": "[ALL] New Years Eve", "start_date": "2020-12-31", "end_date": "2020-12-31"}}} 0x123d570 98 [] false circleci1504710777.zendesk.com map[] map[] <nil> map[]   <nil> <nil> <nil> 0xc0000ac050}
201 Created%                                
```

A sample `holidays.txt` has been provided with a list of holidays for US, UK, Ireland and Japan for 2020. 
