from http.server import HTTPServer, BaseHTTPRequestHandler; import json
class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps([{'id': 1, 'name': 'Item 1', 'status': 'active'}, {'id': 2, 'name': 'Item 2', 'status': 'inactive'}]).encode())
HTTPServer(('', 8888), Handler).serve_forever()
