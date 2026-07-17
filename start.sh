#!/bin/bash

#source ./stop.sh

echo "start server..."

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

        cp gameserver gameserver-$1
        mv gameserver-$1 conf/$1
        echo "start $1..."
        tmp="-log=../../logs/$1 -baseConf=../baseConf.toml"
        for c in `ls conf/$1`
        do
                if [[ $c =~ ^gameserver ]]
                then continue
                fi
	
                array=(`echo $c | tr '.' ' '` )
                #echo $array
	        if [[ ${#array[@]} -ne 2 ]] || [[ $c =~ ^ip  ]]
                then continue
                fi

                tmp="$tmp -$array=$c"
                #echo $tmp
        done
        echo $tmp
        cd conf/$1
        nohup ./gameserver-$1 $tmp >/dev/null 2>&1 &
        exit 0
fi

source ./stop.sh
#cd conf/
for f in `ls conf`
do
	#echo "start $f..."
	if [[ $f =~ ^_ ]] || [[ $f = *.toml ]] || [[ -f $f ]] || [[ $f = *.db ]]
	then 
		echo "pass $f..."
		continue
	fi
	echo "start $f..."
        cp gameserver gameserver-$f
        mv gameserver-$f conf/$f

	tmp="-log=../../logs/$f -baseConf=../baseConf.toml"
	for c in `ls conf/$f`
	do
                if [[ $c =~ ^gameserver|bi ]]
                then continue
                fi
		array=(`echo $c | tr '.' ' '` )
		#echo $array
		if [[ ${#array[@]} -ne 2 ]] || [[ $c =~ ^ip  ]]
		then continue
		fi
		tmp="$tmp -$array=$c"
		#echo $tmp
	done
	#echo $tmp
	#tmp="-module=conf/$f/modules.json -log=logs/$f -conf=conf/$f/config.json -rule=conf/$f/rule.toml"
	#mv gameserver gameserver1
	#nohup ./gameserver -module=conf/$f/modules.json -log=logs/$f -conf=conf/$f/config.json -rule=conf/$f/rule.toml >/dev/null 2>&1 &
	cd conf/$f
	nohup ./gameserver-$f $tmp >/dev/null 2>&1 &
	cd ../../
done

#for f in `ls`
#do
#	if [[ $f =~ ^arena.*(.exe)?$ && -x $f ]]
#	then
#		echo "start $f"
#		nohup ./$f -module=conf/modules.json -log=logs/ -conf=conf/config.json -rule=conf/rule.toml >/dev/null 2>&1 &
#	fi
#done

echo 'finish ...'
