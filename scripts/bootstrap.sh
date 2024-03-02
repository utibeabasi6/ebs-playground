#!/bin/sh

# Create the volume group
vgcreate ebs-pg-vg /dev/xvdf /dev/xvdg /dev/xvdh

# Create logical volumes
sudo lvcreate -l 50%FREE -n lv1  ebs-pg-vg
sudo lvcreate -l 100%FREE -n lv2  ebs-pg-vg

# Create filesystems on logical volumes
sudo mkfs.ext4  /dev/ebs-pg-vg/lv1
sudo mkfs.ext4  /dev/ebs-pg-vg/lv2

# Create mount directory
mkdir -p /home/ubuntu/mounts/lv1
mkdir -p /home/ubuntu/mounts/lv2

# Mount filesystems
sudo su -c "echo '$(sudo blkid /dev/ebs-pg-vg/lv1 | awk '{print $2}') /home/ubuntu/mounts/lv1 ext4 defaults 0 0' >> /etc/fstab"
sudo su -c "echo '$(sudo blkid /dev/ebs-pg-vg/lv2 | awk '{print $2}') /home/ubuntu/mounts/lv2 ext4 defaults 0 0' >> /etc/fstab"
sudo mount -a
