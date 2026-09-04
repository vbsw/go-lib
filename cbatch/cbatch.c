/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <stdbool.h>
#include "cbatch.h"

typedef struct {
	void **data1, **data2;
	const char **names, *err_str;
	int_fast32_t err1, err2, length, index, pass;
} batch_task_params_t;

typedef void (*batch_task_func_t)(batch_task_params_t *params);

void aca_batch_alloc(void ***const data, const int_fast32_t data_len) {
	const size_t size = (size_t)data_len*sizeof(void*);
	void *const data_new = malloc(size);
	memset(data_new, 0, size);
	*data = (void**)data_new;
}

void aca_batch_run(void **const data, int_fast32_t *const err1, int_fast32_t *const err2, int_fast32_t *const err_idx, char **const err_str, const int_fast32_t len, const int_fast32_t passes) {
	batch_task_params_t params = {&data[len], &data[len*2], (const char**)&data[len*3], NULL, 0, 0, len, 0, 0};
	// main
	while (params.pass < passes) {
		// forward
		for (params.index = 0; params.index < len && !params.err1; params.index++) {
			if (data[params.index]) {
				*err_idx = params.index;
				((batch_task_func_t)data[params.index])(&params);
			}
		}
		// backwards
		if (!params.err1 && ++params.pass < passes) {
			params.index = len - 1;
			while (!params.err1) {
				if (data[params.index]) {
					*err_idx = params.index;
					((batch_task_func_t)data[params.index])(&params);
				}
				if (params.index > 0)
					params.index--;
				else
					break;
			}
		}
		if (!params.err1)
			params.pass++;
		else
			break;
	}
	// error handling
	if (params.err1) {
		params.pass = -(params.pass + 1), params.index = len - 1;
		while (true) {
			if (data[params.index]) {
				((batch_task_func_t)data[params.index])(&params);
			}
			if (params.index > 0)
				params.index--;
			else
				break;
		}
		*err1 = params.err1;
		*err2 = params.err2;
		*err_str = (char*)params.err_str;
	}
}

void aca_batch_free(void **const data) {
	free((void*)data);
}
