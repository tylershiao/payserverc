#!/bin/bash

echo 'stop server...'

if [[ $# -gt 0 ]]
then
        echo "stop $1..."
        for f in `ls conf/$1`
        do
        if [[ $f =~ ^gameserver ]]
        then
                echo "stop $f"
                tmp=$(ps -ef | grep $f | grep -v "grep" | awk '{print $2}')
                #echo $tmp
                if [ "$tmp" = "" ];then
                        echo "$f not running"
                else
                        kill $tmp
                fi
                break
        fi
        done
	exit 0
fi

for f in `ls`
do
	if [[ $f =~ ^.*robot(.exe)?|gameserver.*(.exe)?$ && -x $f ]]
	then
		echo "stop $f"
		tmp=$(ps -ef | grep $f | grep -v "grep" | awk '{print $2}')
		#echo $tmp
		if [ "$tmp" = "" ];then
			echo "$f not running"
		else 
			kill $tmp
		fi
	fi
done

#ps -ef | grep "gateway_server\ -port" | grep -v "grep" | awk '{print $2}' | xargs kill -9

echo 'stop ok'
