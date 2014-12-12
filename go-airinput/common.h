#ifndef COMMON_H
#define COMMON_H

struct picture{
  int xres,yres;
  char *buffer;
  struct fb_cmap *colormap;
  char bps,gray;
};

int TakeScreenshot (char* device, struct picture* pict);

#endif
