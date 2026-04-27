/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include <stdlib.h>
#include <stdint.h>
#include <assert.h>
#include "cmodule.h"

void vbsw_cmodule_alloc_buffer(void ***const data, int32_t *const data_len, int32_t *const data_size, const int32_t mod_len_new) {
	const int64_t buffer_len_old = (int64_t)data_len[0]+2;
	const int64_t buffer_len_new = ((int64_t)mod_len_new+1)*2;
	if (buffer_len_old < buffer_len_new) {
		const int64_t buffer_size_new = buffer_len_new*(int64_t)sizeof(void*);
		const int32_t *const buffer_ext1_old = (const int32_t*)(data[0] ? data[0][data_len[0]] : 0);
		const int32_t *const buffer_ext2_old = (const int32_t*)0;
		const int64_t buffer_ext1_size_old = (int64_t)(buffer_ext1_old ? buffer_ext1_old[0] : 0);
		const int64_t buffer_ext2_size_old = (int64_t)0;
		const int64_t data_size_new = buffer_size_new+buffer_ext1_size_old+buffer_ext2_size_old;
		if ((int64_t)(INT32_MAX) >= data_size_new) {
			void **const buffer_new = (void**)malloc((size_t)buffer_size_new);
			if (buffer_new) {
				if (data[0])
					free(data[0]);
				buffer_new[buffer_len_new-2] = (void*)buffer_ext1_old;
				buffer_new[buffer_len_new-1] = (void*)buffer_ext2_old;
				data[0] = buffer_new;
				data_len[0] = buffer_len_new-2;
				data_size[0] = (int32_t)data_size_new;
			} else {
				data_size[0] = 0;
			}
		} else {
			data_size[0] = 0;
		}
	} else {
		data_size[0] = 0;
	}
}

void vbsw_cmodule_proc(void **const data, const int32_t data_len, int32_t *const data_size, const int passes, int32_t *const err_idx, int64_t *const err1, int64_t *const err2, char **const err_str) {
	assert(data);
}

void vbsw_cmodule_free(void **const data, const int32_t data_len) {
	assert(data);
	if (data[data_len])
		free(data[data_len]);
	if (data[data_len+1])
		free(data[data_len+1]);
	free(data);
}
