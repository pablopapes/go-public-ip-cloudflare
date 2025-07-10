# CloudFlare IP Updater
Update Cloudflare DNS with your public IP address using Cloudflare API

## Usage
1. Create a Cloudflare account and get your API key and Zone ID
2. Add a new A record in your Cloudflare account with the name you want to update ( Proxy status: Proxied)
3. Clone this repository
4. edit your .env file with your Cloudflare API key, Zone ID, and DNS Record ID.
EXAMPLE:
```
API_KEY=YOUR_API_KEY
ZONE_ID=YOUR_ZONE_ID
DNS_RECORD_ID=YOUR_DNS_RECORD_ID
```

5. Compile the code
```
go build main.go
```
6. Run the code
```
./main
```
7. You can use a cron job to run the code every X minutes/hours/days
```
*/5 * * * * /path/to/main
```

