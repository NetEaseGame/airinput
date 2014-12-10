/*
 * input.h
 * Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
 *
 * Distributed under terms of the MIT license.
 */

#ifndef INPUT_H
#define INPUT_H

int input_init(char * event);

void tap(int x, int y, int duration_msec);
void drag(int start_x, int start_y, int end_x, int end_y, int num_steps, int msec);
void pinch(int Ax0, int Ay0, int Ax1, int Ay1,
		int Bx0, int By0, int Bx1, int By1, int num_steps, int msec);

#endif /* !INPUT_H */
