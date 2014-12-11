/*
 * screen.h
 * Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
 *
 * Distributed under terms of the MIT license.
 */

#ifndef SCREEN_H
#define SCREEN_H

#include <unistd.h>

void screen_init();

unsigned int width();
unsigned int height();
unsigned int bytespp();
unsigned int offset();


#endif /* !SCREEN_H */
