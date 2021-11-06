#!/bin/sh

dir= # dir with script and list and downladed files
list=list-cinema.lst  # name of file with list 
uid=
usess=
API_KEY=
CHAT_ID=
cd $dir
for i in $(seq 1 10); do
rss=$(curl http://insearch.site/rssdd.xml)
[[ `echo "$rss" | wc -l` -gt 50 ]] && break || sleep 60
done;
[[ `echo "$rss" | wc -l` -le 50 ]] && exit
while read line; do
num=$(echo "$rss" | grep -n ".*$line.*" | grep "1080" | awk -F: '{print$1}')
count=$(echo "$rss" | grep ".*$line.*" | grep "1080" | sed -e "s/<title>//" -e "s/<\/title>//" | tr -d '[:space:]')
for n in $num; do
let link_num=$num+3
link=$(echo "$rss" | sed -n "$link_num"p | sed -e "s/<link>//" -e "s/<\/link>//" | tr -d '[:space:]')
if [[ "`grep "$link" downladed.txt -c`" -ge 1 ]]
then
    break
else
#wget -v --content-disposition --header "Cookie: uid="$uid";usess="$usess"" -- "$link" 
for x in $(seq 1 10); do
echo $link
curl -s --cookie "uid=$uid;usess=$usess" "$link" -O -J && { echo "$link" >> downladed.txt; curl -k --data "text=$count" --data "chat_id=$CHAT_ID" "https://api.telegram.org/bot$API_KEY/sendMessage"; break; } || sleep 60
done;
fi;
done;
done < $list
