#! /bin/bash
function scandir() {
    local cur_dir parent_dir workdir
    workdir=$1
    cd ${workdir}
    if [ ${workdir} = "/" ]
    then
        cur_dir=""
    else
        cur_dir=$(pwd)
    fi

    for dirlist in $(ls ${cur_dir})
    do
        if [[ ${dirlist} = bin || ${dirlist} = *test* || ${dirlist} = db || ${dirlist} = docs || ${dirlist} = pb || ${dirlist} = tools ]] || [[ ! -d ${dirlist} ]]
            then 
            # echo 'passing' ${dirlist}
            continue
        fi
        if test -d ${dirlist};then
            echo 'checking' ${cur_dir}/${dirlist}
            fieldalignment -fix ${cur_dir}/${dirlist}
            cd ${dirlist}
            scandir ${cur_dir}/${dirlist}
            cd ..
        # else
            # c=${cur_dir}/${dirlist}
            # echo $c
	        # if [[ $c = *.go ]]
            # then                 
            #     echo 'checking' ${cur_dir}/${dirlist}
            #     fieldalignment ${cur_dir}/${dirlist}
            # fi
        fi
    done
}

if test -d $1
then
    scandir $1
elif test -f $1
then
    echo "you input a file but not a directory,pls reinput and try again"
    exit 1
else
    echo "the Directory isn't exist which you input,pls input a new one!!"
    exit 1
fi