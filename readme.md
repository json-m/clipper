# clipper

tiny cli for ffmpeg to trim **ReplayBuffer**/**ShadowPlay** recordings into short clips of a target size

will check the output file's size and if larger than desired, will keep reducing quality by re-encoding

made for windows, but probably works on other platforms

### requirements

* NVIDIA GPU with [NVENC support](https://developer.nvidia.com/video-encode-and-decode-gpu-support-matrix-new) (specifically, h264_nvenc)
* ffmpeg built with NVENC is available in your PATH

### how-to

1. setup [obs replay buffer]() or nvidia shadowplay
2. configure clipper.yml with a text editor
3. clip a video, minimal example: `clipper.exe -time=02:14 -dur=12`

this will grab the most recently created file inside inputFolder and create a 12 second long clip, starting at the 02:14

if you wish to specify a specific video file inside inputFolder for clipping, check out the `-file` arg below

### cli arguments

* time - `-time=02:14` - timestamp in video to start at (string)
* duration - `-dur=12` - clip duration in seconds (int, doesn't support 00:00 formatting)
* audio - `-aud=true` - enable/disable audio, specifying this will override config setting (bool)
* file - `-file="Replay 2021-11-09 19-31-21.mkv"` - specify input file (string, filename)

### config opts

* inputFolder - folder to search for input videos (**NO** trailing slash)
* outputFolder - folder to output clips from ffmpeg to (**NO** trailing slash)
* startQuality - decides starting QP level (default: 29)
* targetResolution - target output video res (default: 720p)
* targetFileSize - target output file size in bits (default: 8MB)
* audio - boolean to always include audio or not in output (default: false)
