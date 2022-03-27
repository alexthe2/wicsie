#ffmpeg -r 60 -i out$1/boardgrid%d.png -y -c:v ffv1 -qscale:v 0 simulation.avi -loglevel verbose

ffmpeg -r 24 -i out$1/boardgrid%d.png -y -vcodec png simulation.mov