import os

with open('template.yaml', 'r') as file:
    filedata = file.read()

filedata = filedata.replace('AWS_ACCESS_KEY_ID', os.environ['AWS_ACCESS_KEY_ID'])
filedata = filedata.replace('REGION', os.environ['AWS_DEFAULT_REGION'])
filedata = filedata.replace('AWS_SECRET_ACCESS_KEY', os.environ['AWS_SECRET_ACCESS_KEY'])

with open('file.txt', 'w') as file:
    file.write(filedata)

print(filedata)