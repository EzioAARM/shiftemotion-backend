import os
import base64

with open('template.yaml', 'r') as file:
    filedata = file.read()

spotifykeys = os.environ['SPOTIFY_CLIENT_ID'] + ":" + os.environ['SPOTIFY_SECRET']
spotifybase64 = base64.b64encode(spotifykeys.encode('utf-8'))
filedata = filedata.replace('SPOTIFY_BASE64_KEY', spotifybase64[2: len(spotifybase64) - 1])

with open('file.txt', 'w') as file:
    file.write(filedata)

print(filedata)