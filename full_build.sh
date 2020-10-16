clear
git pull

sudo docker stop am-stats
sudo docker rm am-stats

export PATH=$PATH:/usr/local/go/bin
go build -o build/app

sudo docker build -t cufee/am-stats .
sudo docker run -d --name am-stats -p 6969:4000 cufee/am-stats:latest
