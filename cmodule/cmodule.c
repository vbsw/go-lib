/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include <stdlib.h>
#include <assert.h>
#include "cmodule.h"

#define BUF_SIZE_FUNCDATA(len) (sizeof(void*)*(len*2))
#define BUF_SIZE_LENSORT(len) (sizeof(void*)*(len*2))
#define BUF_SIZE_MODNAMES(len) (sizeof(char*)*(len*80))
#define BUF_SIZE_ERRSTR (sizeof(char)*300)

void vbsw_cmodule_alloc_buffer(void ***const data, size_t *const data_len, void **const data_old, const size_t modules_len) {
	const size_t size = BUF_SIZE_FUNCDATA(modules_len)+BUF_SIZE_LENSORT(modules_len)+BUF_SIZE_MODNAMES(modules_len) + BUF_SIZE_ERRSTR;
	const size_t data_len_new = (size+(sizeof(void*)-1))/sizeof(void*);
	const size_t data_len_size = sizeof(void*)*data_len_new;
	if (data_len[0] < data_len_new && data_len_new < data_len_size) {
		void **const data_new = (void**)malloc(data_len_size);
		if (data_new) {
			data[0] = data_new;
			data_len[0] = data_len_new;
			if (data_old)
				free(data_old);
		}
	}
}

void vbsw_cmodule_proc(void **const data, const size_t modules_len, const int passes, size_t *const err_idx, long long *const err1, long long *const err2, char **const err_str) {
	assert(data);
}

void vbsw_cmodule_rm(void **const data, const size_t modules_len, const int passes, size_t *const err_idx, long long *const err1, long long *const err2, char **const err_str) {
	assert(data);
}
