/*
 * fbshot.c -- FrameBuffer Screen Capture Utility
 * (C)opyright 2002 sfires@sfires.net
 *
 * Originally Written by: Stephan Beyer <PH-Linex@gmx.net>
 * Further changes by: Paul Mundt <lethal@chaoticdreams.org>
 * Rewriten and maintained by: Dariusz Swiderski <sfires@sfires.net>
 * Modular version by: Karol Kuczmarski <karol.kuczmarski@polidea.pl>
 *
 * 	This is a simple program that generates a
 * screenshot of the specified framebuffer device and
 * terminal and writes it to a specified file using
 * the PNG format.
 *
 * See ChangeLog for modifications, CREDITS for credits.
 *
 * fbshot is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either Version 2
 * of the License, or (at your option) any later version.
 *
 * fbshot is distributed in the hope that it will be useful, but
 * WITHOUT ANY  WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public
 * License with fbshot; if not, please write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA
 * 02111-1307 USA
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <fcntl.h>
#include <byteswap.h>
#include <sys/types.h>
#include <asm/types.h>
#include <sys/stat.h>
#include <sys/ioctl.h>
#include <mntent.h>
#include <errno.h>
#include <sys/utsname.h>

#include <sys/vt.h>
#include <linux/fb.h>

#include "common.h"

#define DEFAULT_FB      "/dev/fb0"
#define PACKAGE 	"fbshot"
#define VERSION 	"0.3.1"
#define MAINTAINER_NAME "Dariusz Swiderski"
#define MAINTAINER_ADDR "sfires@sfires.net"

static int waitbfg=0; /* wait before grabbing (for -C )... */


/* some conversion macros */
#define RED565(x)    ((((x) >> (11 )) & 0x1f) << 3)
#define GREEN565(x)  ((((x) >> (5 )) & 0x3f) << 2)
#define BLUE565(x)   ((((x) >> (0)) & 0x1f) << 3)

#define ALPHA1555(x) ((((x) >> (15)) & 0x1 ) << 0)
#define RED1555(x)   ((((x) >> (10)) & 0x1f) << 3)
#define GREEN1555(x) ((((x) >> (5 )) & 0x1f) << 3)
#define BLUE1555(x)  ((((x) >> (0 )) & 0x1f) << 3)

void FatalError(char* err){
  fprintf(stderr,"An error occured: %s %s\nExiting now...\n",err,strerror(errno));
  fflush(stderr);
  exit (1);
}

void Usage(char *binary){
  printf("Usage: %s [-ghi] [-{C|c} vt] [-d dev] [-s n] filename.png\n", binary);
}

void Help(char *binary){
    printf("FBShot - makes screenshots from framebuffer, v%s\n", VERSION);
    printf("\t\tby Dariusz Swiderski <sfires@sfires.net>\n\n");

    Usage(binary);

    printf("\nPossible options:\n");
    printf("\t-C n  \tgrab from console n, for slower framebuffers\n");
    printf("\t-c n  \tgrab from console n\n");
    printf("\t-d dev\tuse framebuffer device dev instead of default\n");
/* not supported as for now
    printf("\t-g    \tsave a grayscaled PNG\n");
 */
    printf("\t-h    \tprint this usage information\n");
    printf("\t-i    \tturns OFF interlacing\n");
    printf("\t-s n  \tsleep n seconds before making screenshot\n");

    printf("\nSend feedback !!!\n");
}

void chvt(int num){
  int fd;
  if(!(fd = open("/dev/console", O_RDWR)))
    FatalError("cannot open /dev/console");
  if (ioctl(fd, VT_ACTIVATE, num))
    FatalError("ioctl VT_ACTIVATE ");
  if (ioctl(fd, VT_WAITACTIVE, num))
    FatalError("ioctl VT_WAITACTIVE");
  close(fd);
  if (waitbfg)
    sleep (3);
}

unsigned int create_bitmask(struct fb_bitfield* bf) {

	return ~(~0u << bf->length) << bf->offset;
}

// Unifies the picture's pixel format to be 32-bit ARGB
void unify(struct picture *pict, struct fb_var_screeninfo *fb_varinfo) {

	__u32 red_mask, green_mask, blue_mask;
	__u32 c;
	__u32 r, g, b;
	__u32* out;
	int i, j = 0, bytes_pp;

	// build masks for extracting colour bits
	red_mask = create_bitmask(&fb_varinfo->red);
	green_mask = create_bitmask(&fb_varinfo->green);
	blue_mask = create_bitmask(&fb_varinfo->blue);

	// go through the image and put the bits in place
	out = (__u32*)malloc(pict->xres * pict->yres * sizeof(__u32));
	bytes_pp = pict->bps >> 3;
	for (i = 0; i < pict->xres * pict->yres * bytes_pp; i += bytes_pp) {

		memcpy (((char*)&c) + (sizeof(__u32) - bytes_pp), pict->buffer + i, bytes_pp);

		// get the colors
		r = ((c & red_mask) >> fb_varinfo->red.offset) & ~(~0u << fb_varinfo->red.length);
		g = ((c & green_mask) >> fb_varinfo->green.offset) & ~(~0u << fb_varinfo->green.length);
		b = ((c & blue_mask) >> fb_varinfo->blue.offset) & ~(~0u << fb_varinfo->blue.length);

		// format the new pixel
		out[j++] = (0xFF << 24) | (b << 16) | (g << 8) | r;
	}

	pict->buffer = (char*)out;
	pict->bps = 32;
}


int read_fb(char *device, int vt_num, struct picture *pict){
  int fd, vt_old, i,j;
  struct fb_fix_screeninfo fb_fixinfo;
  struct fb_var_screeninfo fb_varinfo;
  struct vt_stat vt_info;

  if (vt_num!=-1){
    if ((fd = open("/dev/console", O_RDONLY)) == -1)
      FatalError("could not open /dev/console");
    if (ioctl(fd, VT_GETSTATE, &vt_info))
      FatalError("ioctl VT_GETSTATE");
    close (fd);
    vt_old=vt_info.v_active;
  }

  if(!(fd=open(device, O_RDONLY)))
    FatalError("Couldn't open framebuffer device");

  if (ioctl(fd, FBIOGET_FSCREENINFO, &fb_fixinfo))
    FatalError("ioctl FBIOGET_FSCREENINFO");

  if (ioctl(fd, FBIOGET_VSCREENINFO, &fb_varinfo))
    FatalError("ioctl FBIOGET_VSCREENINFO");

  pict->xres=fb_varinfo.xres;
  pict->yres=fb_varinfo.yres;
  pict->bps=fb_varinfo.bits_per_pixel;
  pict->gray=fb_varinfo.grayscale;

  if(fb_fixinfo.visual==FB_VISUAL_PSEUDOCOLOR){
    pict->colormap=(struct fb_cmap*)malloc(sizeof(struct fb_cmap));
    pict->colormap->red=(__u16*)malloc(sizeof(__u16)*(1<<pict->bps));
    pict->colormap->green=(__u16*)malloc(sizeof(__u16)*(1<<pict->bps));
    pict->colormap->blue=(__u16*)malloc(sizeof(__u16)*(1<<pict->bps));
    pict->colormap->transp=(__u16*)malloc(sizeof(__u16)*(1<<pict->bps));
    pict->colormap->start=0;
    pict->colormap->len=1<<pict->bps;
    if (ioctl(fd, FBIOGETCMAP, pict->colormap))
      FatalError("ioctl FBIOGETCMAP");
  }
  if (vt_num!=-1)
    chvt(vt_old);

  switch(pict->bps){
  case 15:
    i=2;
    break;
  default:
    i=pict->bps>>3;
  }

  if(!(pict->buffer=malloc(pict->xres*pict->yres*i)))
    FatalError("couldnt malloc");

  fprintf(stdout, "Framebuffer %s is %i bytes.\n", device,
                    (fb_varinfo.xres * fb_varinfo.yres * i));
  fprintf(stdout, "Grabbing %ix%i ... \n", fb_varinfo.xres, fb_varinfo.yres);

//#ifdef DEBUG
///* Output some more information bout actual graphics mode
// */
//  fprintf(stdout, "%ix%i [%i,%i] %ibps %igr\n",
//  	fb_varinfo.xres_virtual, fb_varinfo.yres_virtual,
//  	fb_varinfo.xoffset, fb_varinfo.yoffset,
//  	fb_varinfo.bits_per_pixel, fb_varinfo.grayscale);
//    fprintf(stdout, "FIX: card:%s mem:0x%.8X mem_len:%d visual:%i type:%i type_aux:%i line_len:%i accel:%i\n",
//  fb_fixinfo.id,fb_fixinfo.smem_start,fb_fixinfo.smem_len,fb_fixinfo.visual,
//  fb_fixinfo.type,fb_fixinfo.type_aux,fb_fixinfo.line_length,fb_fixinfo.accel);
//#endif

  fflush(stdout);
  if (vt_num!=-1)
    chvt(vt_num);

  j= (read(fd, pict->buffer, ((pict->xres * pict->yres) * i) )!=
  	(pict->xres * pict->yres *i ));
//#ifdef DEBUG
//  printf("to read:%i read:%i\n",(pict->xres* pict->yres * i), j);
//#endif
  if (vt_num!=-1)
    chvt(vt_old);

//  if(j)
//    FatalError("couldn't read the framebuffer");
//  else
//    fprintf(stdout,"done.\n");
  close (fd);

  unify(pict, &fb_varinfo);
  return 0;
}


void convert8to32(struct picture *pict){
  int i;
  int j=0;
  __u8 c;
  char *out=(char*)malloc(pict->xres*pict->yres*4);
  for (i=0; i<pict->xres*pict->yres; i++)
  {
    c = ((__u8*)(pict->buffer))[i];
    out[j++]=(char)(pict->colormap->red[c]);
    out[j++]=(char)(pict->colormap->green[c]);
    out[j++]=(char)(pict->colormap->blue[c]);
    out[j++]=(char)(pict->colormap->transp[c]);
  }
  free(pict->buffer);
  pict->buffer=out;
}

void convert1555to32(struct picture *pict){
  int i;
  int j=0;
  __u16 t,c;
  char *out=(char*)malloc(pict->xres*pict->yres*4);
  for (i=0; i<pict->xres*pict->yres; i++)
  {
    c = ( (__u16*)(pict->buffer))[i];
    out[j++]=(char)RED1555(c);
    out[j++]=(char)GREEN1555(c);
    out[j++]=(char)BLUE1555(c);
    out[j++]=(char)ALPHA1555(c);
  }
  free(pict->buffer);
  pict->buffer=out;
}

void convert565to32(struct picture *pict){ // ARGB_8888
  int i;
  int j=0;
  __u16 t,c;
  char *out=(char*)malloc(pict->xres*pict->yres*4);
  for (i=0; i<pict->xres*pict->yres; i++)
  {
    c = ( (__u16*)(pict->buffer))[i];
    out[j++]=(char)0xff;
    out[j++]=(char)RED565(c);
    out[j++]=(char)GREEN565(c);
    out[j++]=(char)BLUE565(c);

  }
  free(pict->buffer);
  pict->buffer=out;
}


static int Write(int fd, char* buf, int c)
{
	int i = 0, r;
	while (i < c)
	{
		r = write(fd, (const void*)(buf + i), c - i);
		i += r;
	}
	return i;
}


#if 0
static int Write_RAW(struct picture* pict, char* filename)
{
  int fd = open(filename, O_CREAT | O_WRONLY | O_SYNC, 0777);
  if (fd < 0)	return -1;

  // convert picture to 32-bit format
  printf ("BPS: %d\n", pict->bps);
  switch (pict->bps)
  {
  	case 8:
		convert8to32 (pict);
		break;
	case 15:
		convert1555to32 (pict);
		break;
	case 16:
		convert565to32 (pict);
		break;
	case 32:
		break;
	default:
		return -1;
  }

  Write (fd, pict->buffer, pict->xres*pict->yres*4);
  close (fd);

  return 0;
}
#endif


static char optstring[] = "hiC:c:d:s:";
static struct option long_options[] = {
        {"slowcon", 1, 0, 'C'},
        {"console", 1, 0, 'c'},
        {"device", 1, 0, 'd'},
        {"help", 0, 0, 'h'},
        {"noint", 0, 0, 'i'},
        {"sleep", 1, 0, 's'},
        {0, 0, 0, 0}
        };


int TakeScreenshot (char* device, struct picture* pict)
{
	int vt_num = -1;
	return read_fb(device, vt_num, pict);
}
