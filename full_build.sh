clear
git pull
export PATH=$PATH:/usr/local/go/bin
go build -o build/app

sudo docker stop am-stats
sudo docker rm am-stats

sudo docker build -t cufee/am-stats .
sudo docker run -d --restart unless-stopped --name am-stats -p 6969:4000 cufee/am-stats:latest
