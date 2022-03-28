mkdir /out/compiled/$1
ffmpeg -r 60 -i out/raw/boardgrid%d.png -y -vcodec png out/compiled/$1/simulation_$2.mov