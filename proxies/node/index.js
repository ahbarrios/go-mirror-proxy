const http = require('http');
const { URL } = require('url');
const args = process.argv.slice(2);

let port = 8080;
let host = 'localhost';
for (let i = 0; i < args.length; i++) {
    let arg = args[i];
    if (arg === '--addr' && i + 1 < args.length) {
        const url = new URL(args[i + 1]);
        host = url.hostname;
        port = url.port;
    }
}

let server = http.createServer();

// proxy
server.on('request', (req, res) => {
	console.log("Received request for ", req.url)
    const { hostname, port, pathname } = new URL(req.url);
    http.request({
        port,
        host: hostname,
        path: pathname,
        method: req.method,
        headers: req.headers
    }, (response) => {
        res.writeHead(response.statusCode, response.headers);
        response.pipe(res);
    }).on('error', (err) => {
        console.error(err);
        res.writeHead(500);
        res.end();
    }).end();
})

server.listen(port, host, () => {
    console.log('server is running at 3000');
})
