rsync -azv --exclude=node_modules  --exclude=web2 --exclude=*.exe  --exclude '.git' --delete . ubuntu@node4:/home/ubuntu/go/src/github.com/sfi2k7/picoweb