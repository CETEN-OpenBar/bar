FROM ghcr.io/cedws/concrete-ubuntu:builder AS build
LABEL project=bar-borne

WORKDIR /root

RUN wget -O rootfs.tar.gz https://cdimage.ubuntu.com/ubuntu-base/releases/22.04/release/ubuntu-base-22.04-base-amd64.tar.gz

RUN mkdir rootfs/
RUN tar -C rootfs/ -xf rootfs.tar.gz

COPY ./fs/etc rootfs/etc

RUN fakeroot fakechroot chroot rootfs/ \
    apt-get update

RUN chroot rootfs/ \
    apt-get install -yq --no-install-recommends \
        linux-image-generic initramfs-tools parted overlayroot systemd systemd-sysv dbus sudo

RUN fakeroot fakechroot chroot rootfs/ \
    apt-get clean

RUN chroot rootfs/ \
    useradd -m ubuntu -G sudo -s /bin/bash
RUN chroot rootfs/ \
    bash -c 'echo "ubuntu:ubuntu" | chpasswd'

# Only needed during install to prevent kernel/initrd symlinks from being created
RUN chroot rootfs/ \
    rm /etc/kernel-img.conf

# Usefull stuff for our bar purposes :
# Install npm, git, xorg and chromium
# clone https://github.com/yyewolf/bar
# go to bar/frontend
# setup config.json with the right values :
#     echo {"api":"http://localhost:8080","apiws":"ws://localhost:8080","local_token": "jesuisleplusfort"} > bar/frontend/static/config.json
# do npm install and setup a service that does node build/index.js at boot in bar/frontend
# made it so x server starts at boot and chromium is in kiosk mode as soon as x server starts
# made it so chromium starts at boot and goes to localhost:8080

COPY bar.service .
COPY override.conf .
COPY .bashrc .
COPY .xinitrc .

RUN cp bar.service rootfs/etc/systemd/system/bar.service
RUN mkdir rootfs/etc/systemd/system/getty@tty1.service.d
RUN cp override.conf rootfs/etc/systemd/system/getty@tty1.service.d/override.conf
RUN cp .bashrc rootfs/home/ubuntu/.bashrc
RUN cp .xinitrc rootfs/home/ubuntu/.xinitrc

RUN chroot rootfs/ \
    apt-get install -y --no-install-recommends \
        curl git xorg chromium-browser unclutter

# download nvm install script
RUN apt update && apt install curl bzip2 -y
RUN curl https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh -o install_nvm.sh
RUN mv install_nvm.sh rootfs/home/ubuntu/install_nvm.sh
RUN curl -L -o firefox.tar.bz2 "https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=en-US" 
RUN tar xjf firefox.tar.bz2 -C rootfs/usr/lib
RUN chmod +x rootfs/usr/lib/firefox/firefox
RUN ln -s /usr/lib/firefox/firefox /usr/bin/firefox

RUN chroot rootfs/ /bin/bash <<EOF
cd ~/ && \
git clone https://github.com/nvm-sh/nvm.git .nvm && \
cd .nvm && \
. ./nvm.sh && \
cd /home/ubuntu && \
cat install_nvm.sh | /bin/bash && \
nvm install --lts && \
nvm use --lts && \
git clone https://github.com/yyewolf/bar && \
cd bar/frontend && \
npm ci && \
npm run build && \
echo {"api":"http://localhost:8080","apiws":"ws://localhost:8080","local_token": "jesuisleplusfort"} > static/config.json  && \
systemctl enable bar.service
EOF
    

COPY ./utils/sectors.sh .

# Create empty file for containing the partitions
RUN fallocate -l 4GiB image.img

# Set up partition table
RUN parted -s -- image.img \
    mktable gpt \
    mkpart primary 0% 100MiB \
    mkpart primary 100MiB -16MiB \
    mkpart primary -16MiB 100% \
    set 1 esp on

RUN mkdir -p esp/EFI/BOOT esp/EFI/ubuntu esp/grub

COPY ./fs/boot/grub/grub.cfg esp/grub/
COPY ./fs/boot/EFI/ubuntu/grub.cfg ./esp/EFI/ubuntu/
RUN cp -r /usr/lib/grub/x86_64-efi esp/grub/x86_64-efi

# Make unsigned GRUB EFI image
RUN grub-mkimage -o esp/EFI/BOOT/BOOTX64.efi -O x86_64-efi -p "(,gpt1)/grub" fat part_gpt

RUN cp /usr/lib/shim/BOOTX64.CSV esp/EFI/ubuntu/BOOTX64.CSV
RUN cp /usr/lib/shim/shimx64.efi.signed esp/EFI/ubuntu/shimx64.efi
RUN cp /usr/lib/shim/mmx64.efi esp/EFI/ubuntu/mmx64.efi
RUN cp /usr/lib/grub/x86_64-efi-signed/grubx64.efi.signed esp/EFI/ubuntu/grubx64.efi

# Create new FAT image for the ESP
RUN mkfs.fat -C esp.img \
    $(./sectors.sh size image.img 1)

# Copy the contents of esp/ into the FAT esp.img
# Then dd it into the main image without truncating
RUN mcopy -s -i esp.img esp/* ::

# Make new EXT4 image
# Specify the LABEL should be OVERLAY which makes mounting it easier later on
# Create it with the calculated correct size, bearing in mind default block size of 1024
RUN mkfs.ext4 -L OVERLAY overlay.img \
    $(( $(./sectors.sh size image.img 3) / 2 ))

# Combine the individual partitions into the main image
# For each partition, seek to the sector it should start at
# conv=notrunc ensures that no truncation of the image occurs after each command
RUN dd if=esp.img of=image.img status=progress conv=notrunc \
    seek=$(./sectors.sh start image.img 1)
RUN dd if=overlay.img of=image.img status=progress conv=notrunc \
    seek=$(./sectors.sh start image.img 3)

# Pack the rootfs into a squash image
# GRUB is able to load the kernel and initramfs inside on boot
RUN mkdir -p squash
RUN mksquashfs rootfs/ squash/rootfs.squash

# Make new EXT4 image initialised with contents of squash/
# Specify the UUID we want to create it with which matches up with our GRUB config
# Create it with the calculate correct size, bearing in mind default block size of 1024
RUN mkfs.ext4 -d squash/ squash.img -U a185a4b8-f920-4c5d-8b19-8d5a9e49024e \
    $(( $(./sectors.sh size image.img 2) / 2 ))

RUN dd if=squash.img of=image.img status=progress conv=notrunc \
    seek=$(./sectors.sh start image.img 2)

FROM scratch AS export
LABEL project=bar-borne
COPY --from=build /root/squash/rootfs.squash /root/image.img /
