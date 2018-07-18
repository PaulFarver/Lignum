#!/usr/bin/env python3

from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse
import json
import csv
import random

class Tree:
    def __init__(self, genus, species):
        self.genus = genus
        self.species = species

trees = []
with open('global_tree_search_trees_1_2.csv', 'rt') as csvfile:
    reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in reader:
        t = row[0].split()
        trees.append(Tree(t[0], t[1]))

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps(trees[random.randint(0, len(trees) - 1)].__dict__).encode())
        return

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 80), RequestHandler)
    print('Starting server at http://0.0.0.0:80')
    server.serve_forever()