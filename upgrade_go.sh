sudo apt remove --purge golang-go
sudo apt remove --purge gccgo-go
sudo rm -rf /usr/local/go

# Download latest Go version (1.22.0 as of February 2024)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Set up environment variables (add to ~/.profile or ~/.bashrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.profile
source ~/.profile