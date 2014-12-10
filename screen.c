/*
 * screen.c
 * Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
 *
 * Distributed under terms of the MIT license.
 */

#include "screen.h"

#include <fcntl.h>
#include <sys/ioctl.h>
#include <linux/fb.h>

struct fb_var_screeninfo vinfo;

void screen_init(){
	int fd = open("/dev/graphics/fb0", O_RDONLY);
	ioctl(fd, FBIOGET_VSCREENINFO, &vinfo);
	close(fd);
}

unsigned int width(){
	return vinfo.xres;
}

unsigned int height(){
	return vinfo.yres;
}

unsigned int bytespp(){
	return vinfo.bits_per_pixel/8;
}

unsigned int offset(){
	return (vinfo.xoffset + vinfo.yoffset*vinfo.xres) * bytespp();
}
