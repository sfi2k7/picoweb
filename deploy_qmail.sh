#rsync -azv --exclude 'jobs'  --exclude '.git' --delete . root@bmysql:/root/gowork/src/github.com/faisal/blueq
rsync -azv --exclude=blueregister  --exclude 'jobs' --exclude=*.exe  --exclude '.git' --delete . ubuntu@qmailserver:/home/ubuntu/go/src/github.com/sfi2k7/picoweb