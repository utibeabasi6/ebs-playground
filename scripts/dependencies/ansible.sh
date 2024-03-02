#!/bin/bash

apt install python3-pip
python3 -m pip install --user ansible  --break-system-packages
 echo "export PATH=$PATH:$HOME/.local/bin" >> ~/.bashrc
 source ~/.bashrc