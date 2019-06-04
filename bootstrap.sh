# Install Node 12
curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -
sudo apt-get install -y nodejs build-essential

# Install Vue CLI
sudo npm i -g @vue/cli

# Install FFmpeg
sudo add-apt-repository -y ppa:jonathonf/ffmpeg-4
sudo apt-get update
sudo apt-get install -y ffmpeg
